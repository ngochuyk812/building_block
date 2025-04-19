package helpers

import (
	"context"

	auth_context "github.com/ngochuyk812/building_block/pkg/auth"
)

type contextKey string

const (
	AuthContextKey  contextKey = "auth_context"
	TokenContextKey contextKey = "token_context"
)

func AuthContext(ctx context.Context) (*auth_context.AuthContext, bool) {
	val, ok := ctx.Value(AuthContextKey).(*auth_context.AuthContext)
	return val, ok
}

func TokenContext(ctx context.Context) string {
	val, ok := ctx.Value(TokenContextKey).(string)
	if !ok {
		return ""
	}
	return val
}

func SetTokenContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenContextKey, token)
}

func FromContext(ctx context.Context, key contextKey) (*auth_context.AuthContext, bool) {
	val, ok := ctx.Value(key).(*auth_context.AuthContext)
	return val, ok
}

func NewContext(ctx context.Context, key contextKey, authCtx *auth_context.AuthContext) context.Context {
	return context.WithValue(ctx, key, authCtx)
}
