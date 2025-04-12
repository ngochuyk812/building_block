package cache

import (
	"context"
	"time"
)

type ICache interface {
	Get(ctx context.Context, key string, scan interface{}) error
	Set(ctx context.Context, key string, value interface{}, expired time.Duration) error
	Del(ctx context.Context, key string) error
	Gets(ctx context.Context, keys []string, scan map[string]interface{}) error
	Dels(ctx context.Context, keys ...string) error
	Sets(ctx context.Context, data map[string]interface{}, expired time.Duration) error
	WithPrefix(key ...string) (outputKey string)
}
