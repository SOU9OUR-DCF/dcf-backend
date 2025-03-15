package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/adapters/config"
)

type Connection struct {
	Client *redis.Client
}

func NewConnection(cfg *config.Config) (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Connection{Client: client}, nil
}

func (c *Connection) Close() {
	c.Client.Close()
}