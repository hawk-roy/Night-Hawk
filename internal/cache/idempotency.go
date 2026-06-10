package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	processingTTL   = 10 * time.Minute
	successTTL      = 24 * time.Hour
	processingValue = "PROCESSING"
	successPrefix   = "SUCCESS:"
)

type IdempotencyManager struct {
	client *redis.Client
}

func NewIdempotencyManager(client *redis.Client) *IdempotencyManager {
	return &IdempotencyManager{client: client}
}

func (m *IdempotencyManager) key(userID int64, idempotencyKey string) string {
	return fmt.Sprintf("order:idempotency:%d:%s", userID, idempotencyKey)
}

func (m *IdempotencyManager) Acquire(ctx context.Context, userID int64, idempotencyKey string) (bool, error) {
	ok, err := m.client.SetNX(ctx, m.key(userID, idempotencyKey), processingValue, processingTTL).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (m *IdempotencyManager) MarkSuccess(ctx context.Context, userID int64, idempotencyKey string, orderNo string) error {
	return m.client.Set(ctx, m.key(userID, idempotencyKey), successPrefix+orderNo, successTTL).Err()
}

func (m *IdempotencyManager) Release(ctx context.Context, userID int64, idempotencyKey string) error {
	return m.client.Del(ctx, m.key(userID, idempotencyKey)).Err()
}
