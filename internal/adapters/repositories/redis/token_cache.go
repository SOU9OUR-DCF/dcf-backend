package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
)

type tokenCache struct {
	conn *Connection
}

func NewTokenCache(conn *Connection) ports.TokenCache {
	return &tokenCache{
		conn: conn,
	}
}

func (c *tokenCache) StoreToken(ctx context.Context, userID uuid.UUID, token string, expiration time.Duration) error {
	// Store by user ID
	userKey := fmt.Sprintf("user:%s:token", userID.String())
	if err := c.conn.Client.Set(ctx, userKey, token, expiration).Err(); err != nil {
		return err
	}

	// Store by token for validation
	tokenKey := fmt.Sprintf("token:%s", token)
	return c.conn.Client.Set(ctx, tokenKey, userID.String(), expiration).Err()
}

func (c *tokenCache) GetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	key := fmt.Sprintf("user:%s:token", userID.String())
	token, err := c.conn.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *tokenCache) InvalidateToken(ctx context.Context, userID uuid.UUID) error {
	// Get the token first
	token, err := c.GetToken(ctx, userID)
	if err != nil {
		return err
	}

	// Delete by user ID
	userKey := fmt.Sprintf("user:%s:token", userID.String())
	if err := c.conn.Client.Del(ctx, userKey).Err(); err != nil {
		return err
	}

	// Delete by token
	tokenKey := fmt.Sprintf("token:%s", token)
	return c.conn.Client.Del(ctx, tokenKey).Err()
}

func (c *tokenCache) TokenExists(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("token:%s", token)
	exists, err := c.conn.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
