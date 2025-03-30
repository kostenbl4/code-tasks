package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, value any) error
	Set(ctx context.Context, key string, value any, TTL time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}