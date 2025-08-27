package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

type Adapter struct {
	client  *redis.Client
	context context.Context
}

func New(ctx context.Context, config Config) (*Adapter, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Adapter{client: rdb, context: ctx}, nil
}

func (a *Adapter) Client() *redis.Client {
	return a.client
}

func (a *Adapter) Context() context.Context {
	return a.context
}

func (a *Adapter) Close() error {
	if a == nil || a.client == nil {
		return nil
	}

	return a.client.Close()
}
