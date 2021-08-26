package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient interacts with Redis.
type Client interface {
	Get(ctx context.Context, key string) (*string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Close() error
}

// Config specifies how the client should be configured.
type Config struct {
	URI string
}

// Redis is a standalone client for Redis.
type Redis struct {
	c *redis.Client
}

// Get retrieves a key.
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

// Set sets a key with optional expiration.
//
// Expiration of 0 means that the key will not have an expiration.
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.c.Set(ctx, key, value, expiration).Err()
}

// TTL returns the key's expiration as a duration.
//
// -2 means that the key does not exist and -1 means that the key exists but does not have an expiration set.
func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.c.TTL(ctx, key).Result()
}

// Close shuts down the connection to Redis.
func (r *Redis) Close() error {
	return r.c.Close()
}

// New creates a new Redis standalone client.
func New(config Config) (*Redis, error) {
	options, err := redis.ParseURL(config.URI)
	if err != nil {
		return nil, err
	}

	return &Redis{c: redis.NewClient(options)}, nil
}
