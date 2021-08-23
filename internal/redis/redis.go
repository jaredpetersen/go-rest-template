package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient interacts with Redis
type Client interface {
	Get(ctx context.Context, key string) (*string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Close() error
}

type Config struct {
	URI string
}

type Redis struct {
	c *redis.Client
}

func (r *Redis) Get(ctx context.Context, key string) (*string, error) {
	val, err := r.c.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.c.Set(ctx, key, value, 0).Err()
}

func (r *Redis) Close() error {
	return r.c.Close()
}

func New(config Config) (*Redis, error) {
	options, err := redis.ParseURL(config.URI)
	if err != nil {
		return nil, err
	}

	return &Redis{c: redis.NewClient(options)}, nil
}
