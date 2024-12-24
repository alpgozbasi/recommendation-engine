package cache

import (
	"context"
	"time"

	"github.com/alpgozbasi/recommendation-engine/internal/config"
	"github.com/alpgozbasi/recommendation-engine/internal/util"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg *config.AppConfig) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// test connection with a ping
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		util.Logger.Error().Err(err).Msg("failed to connect to Redis")
		return nil, err
	}

	util.Logger.Info().Msg("connected to Redis successfully")
	return &RedisCache{client: rdb}, nil
}

func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}
