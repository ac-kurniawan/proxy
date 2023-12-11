package core

import (
	"context"
	"time"

	"github.com/ac-kurniawan/proxy/core/model"
)

type IProxyService interface {
	RenewExternalToken(ctx context.Context)
	Process(ctx context.Context, method, path, token string, payload map[string]interface{}) (map[string]interface{}, error)
	Register(ctx context.Context, username, password, firstName, lastName, telephone, address, city, province, country string, image []byte) (map[string]interface{}, error)
}

type ProxyService struct {
	Repository IRepository
	Util       IUtil
}

// RenewExternalToken implements IProxyService.
func (p *ProxyService) RenewExternalToken(ctx context.Context) {
	ctx, span := p.Util.StartTrace(ctx, "[SERVICE] RenewExternalToken")
	defer p.Util.EndTrace(span)
	timeNow := time.Now()
	almostExp := timeNow.Add(time.Hour * 1)
	users, err := p.Repository.GetAlmostExpToken(ctx, timeNow, almostExp)
	if err != nil {
		p.Util.TraceError(span, err)
		p.Util.LogError(ctx, err)
		return
	}

	for _, val := range users {
		// TODO: change this function to more atomic
		token := val.ExternalToken
		res, err := p.Repository.Call(ctx, "POST", "/api/api/token/refresh", &token, map[string]interface{}{
			"refresh": val.ExternalRefreshToken,
		})
		if err != nil {
			p.Util.TraceError(span, err)
			p.Util.LogError(ctx, err)
		}
		val.ExternalToken = res["access_token"].(string)
		val.ExternalRefreshToken = res["refresh_token"].(string)
		externalExpToken := p.Util.GetExpFromToken(res["access_token"].(string))
		val.ExternalTokenExpired = time.Unix(externalExpToken, 0)
		p.Repository.Save(ctx, val)
	}
}

// Register implements IProxyService.
func (p *ProxyService) Register(ctx context.Context, username, password, firstName, lastName, telephone, address, city, province, country string, image []byte) (map[string]interface{}, error) {
	ctx, span := p.Util.StartTrace(ctx, "[SERVICE] Register")
	defer p.Util.EndTrace(span)
	response, err := p.Repository.Call(ctx, "POST", "/api/register", nil, map[string]interface{}{
		"username":      username,
		"password":      password,
		"first_name":    firstName,
		"last_name":     lastName,
		"telephone":     lastName,
		"address":       address,
		"city":          city,
		"province":      province,
		"country":       country,
		"profile_image": image,
	})
	if err != nil {
		p.Util.TraceError(span, err)
		p.Util.LogError(ctx, err)
		return nil, err
	}
	timeNow := time.Now()
	exp := timeNow.Add(time.Hour * 24).Unix()             // access token exp in 1 day
	refrestExp := timeNow.Add(time.Hour * 24 * 30).Unix() // refresh token exp in 30 days
	token := p.Util.GenerateJWT(map[string]interface{}{
		"username": username,
	}, exp)
	refreshToken := p.Util.GenerateJWT(map[string]interface{}{
		"username": username,
	}, refrestExp)

	externalExpToken := p.Util.GetExpFromToken(response["access_token"].(string))
	userModel := model.UserModel{
		Username:             username,
		Password:             password,
		ExternalToken:        response["access_token"].(string),
		ExternalRefreshToken: response["refresh_token"].(string),
		ExternalTokenExpired: time.Unix(externalExpToken, 0),
		InternalRefreshToken: refreshToken,
	}

	err = p.Repository.Save(ctx, userModel)
	if err != nil {
		p.Util.TraceError(span, err)
		p.Util.LogError(ctx, err)
		return nil, err
	}

	return map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
	}, nil
}

// Process implements IProxyService.
func (p *ProxyService) Process(ctx context.Context, method, path, username string, payload map[string]interface{}) (map[string]interface{}, error) {
	ctx, span := p.Util.StartTrace(ctx, "[SERVICE] Process")
	defer p.Util.EndTrace(span)
	user, err := p.Repository.FindByUsername(ctx, username)
	if err != nil {
		p.Util.TraceError(span, err)
		p.Util.LogError(ctx, err)
		return nil, err
	}
	token := user.ExternalToken
	return p.Repository.Call(ctx, method, path, &token, payload)
}

func NewProxyService(module ProxyService) IProxyService {
	return &module
}
