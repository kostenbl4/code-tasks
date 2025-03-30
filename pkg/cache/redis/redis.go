package rediscache

import (
	"code-tasks/pkg/cache"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host         string        `env:"CACHE_HOST" yaml:"host" env-required:"true"`
	Port         int           `env:"CACHE_PORT" yaml:"port" env-required:"true"`
	Password     string        `env:"CACHE_PASSWORD" yaml:"password" env-required:"true"`
	TTL          time.Duration `env:"CACHE_TTL" yaml:"TTL" env-default:"30min"`
	WriteTimeout time.Duration `env:"CACHE_WRITE_TIMEOUT" yaml:"write_timeout" env-default:"3s"`
	ReadTimeout  time.Duration `env:"CACHE_READ_TIMEOUT" yaml:"read_timeout" env-default:"2s"`
}

type redisCache struct {
	client *redis.Client
}

func NewRedis(cfg RedisConfig) (cache.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return &redisCache{client: client}, nil
}

func (r *redisCache) Get(ctx context.Context, key string, value any) error {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("failed to get data from redis: %w", err)
	}
	if err := json.Unmarshal(val, value); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return nil
}

func (r *redisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := r.client.Set(ctx, key, val, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set data in redis: %w", err)
	}
	return nil
}

func (r *redisCache) Delete(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete data from redis: %w", err)
	}
	return nil
}
