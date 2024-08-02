package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type BaseRedis struct {
	redis *redis.Client
}

func (r *BaseRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.redis.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *BaseRedis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *BaseRedis) Delete(ctx context.Context, key string) error {
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
