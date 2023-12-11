package core

import "context"

type IUtil interface {
	GenerateJWT(attribute map[string]interface{}, exp int64) string
	GetExpFromToken(token string) int64
	EncryptPassword(password string) string
	StartTrace(ctx context.Context, name string) (context.Context, any)
	EndTrace(span any)
	TraceError(span any, error error)
	LogInfo(ctx context.Context, log string)
	LogError(ctx context.Context, err error)
}
