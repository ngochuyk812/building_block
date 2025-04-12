package helpers

import (
	"context"

	auth_context "github.com/ngochuyk812/building_block/pkg/auth"
)

type contextKey string

const (
	AuthContextKey contextKey = "auth_context"
)

func AuthContext(ctx context.Context) (*auth_context.AuthContext, bool) {
	val, ok := ctx.Value(AuthContextKey).(*auth_context.AuthContext)
	return val, ok
}

func FromContext(ctx context.Context, key contextKey) (*auth_context.AuthContext, bool) {
	val, ok := ctx.Value(key).(*auth_context.AuthContext)
	return val, ok
}

func NewContext(ctx context.Context, key contextKey, authCtx *auth_context.AuthContext) context.Context {
	return context.WithValue(ctx, key, authCtx)
}
