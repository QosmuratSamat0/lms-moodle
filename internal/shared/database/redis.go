// internal/shared/database/redis.go
package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ap1-final-mini-moodle/internal/shared/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

// NullRedisClient is a no-op Redis client used when Redis is unavailable
type NullRedisClient struct{}

func (n *NullRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil // no-op
}

func (n *NullRedisClient) Get(ctx context.Context, key string) (string, error) {
	return "", nil // no-op, returns empty string
}

func (n *NullRedisClient) Del(ctx context.Context, key string) error {
	return nil // no-op
}

func (n *NullRedisClient) Close() error {
	return nil // no-op
}

func (n *NullRedisClient) Health(ctx context.Context) error {
	return nil // no-op
}

func NewRedis(cfg config.RedisURL) (*RedisClient, error) {
	opts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	}

	if cfg.TLSEnabled {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in redis: %w", err)
	}
	return nil

}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key does not exist
		}
		return "", fmt.Errorf("failed to get key from redis: %w", err)
	}
	return val, nil
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key from redis: %w", err)
	}
	return nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

func (r *RedisClient) Health(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
