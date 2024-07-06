package redisconn

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"wildwest/pkg/settings"
)

func NewRedisClient(cfg *settings.Config) (*redis.Client, error) {
	ctx := context.Background()
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	}
	client := redis.NewClient(options)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return client, nil
}
