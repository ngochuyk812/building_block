package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	redis *redis.Client
}

var _ ICache = (*RedisCache)(nil)

func NewRedisCache(connectString, pass string) (ICache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     connectString,
		Password: pass,
		DB:       0,
	})
	fmt.Println("Inject Redis...")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Redis: %v", err)
	}
	fmt.Println("Inject Redis susccessfull")
	return &RedisCache{redis: rdb}, nil
}

func (c RedisCache) Get(ctx context.Context, key string, scan interface{}) error {
	data, err := c.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, scan)
}

func (c RedisCache) Set(ctx context.Context, key string, value interface{}, expired time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, data, expired).Err()
}

func (c RedisCache) Del(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}

func (c RedisCache) Gets(ctx context.Context, keys []string, scan map[string]interface{}) error {
	results, err := c.redis.MGet(ctx, keys...).Result()
	if err != nil {
		return err
	}

	for i, raw := range results {
		if raw == nil {
			continue
		}
		bytesVal, ok := raw.(string)
		if !ok {
			continue
		}
		err := json.Unmarshal([]byte(bytesVal), scan[keys[i]])
		if err != nil {
			return err
		}
	}
	return nil
}

func (c RedisCache) Dels(ctx context.Context, keys ...string) error {
	return c.redis.Del(ctx, keys...).Err()
}

func (c RedisCache) Sets(ctx context.Context, data map[string]interface{}, expired time.Duration) error {
	pipe := c.redis.Pipeline()
	for k, v := range data {
		val, err := json.Marshal(v)
		if err != nil {
			return err
		}
		pipe.Set(ctx, k, val, expired)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (c RedisCache) WithPrefix(prefix ...string) (out string) {
	return strings.Join(prefix, ">")
}
