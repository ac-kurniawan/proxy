package core

import (
	"context"
	"time"

	"github.com/ac-kurniawan/proxy/core/model"
)

type IProxyDB interface {
	Save(ctx context.Context, userModel model.UserModel) error
	FindByUsername(ctx context.Context, username string) (*model.UserModel, error)
	GetAlmostExpToken(ctx context.Context, from, to time.Time) ([]model.UserModel, error)
}

type IDataSource interface {
	Call(ctx context.Context, method, path string, token *string, payload map[string]interface{}) (map[string]interface{}, error)
}

type IRepository interface {
	IDataSource
	IProxyDB
}
