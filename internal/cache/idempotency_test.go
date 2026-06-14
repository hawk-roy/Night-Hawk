package cache

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestIdempotencyManager(t *testing.T) (*IdempotencyManager, *miniredis.Miniredis, *redis.Client) {
	t.Helper()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	manager := NewIdempotencyManager(client)

	t.Cleanup(func() {
		_ = client.Close()
		s.Close()
	})

	return manager, s, client
}

func TestIdempotencyAcquireFirstTime(t *testing.T) {
	manager, _, _ := newTestIdempotencyManager(t)
	ctx := context.Background()

	acquired, err := manager.Acquire(ctx, 1, "key-1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !acquired {
		t.Fatal("expected acquired=true on first acquire")
	}
}

func TestIdempotencyAcquireDuplicate(t *testing.T) {
	manager, _, _ := newTestIdempotencyManager(t)
	ctx := context.Background()

	acquired, err := manager.Acquire(ctx, 1, "key-1")
	if err != nil {
		t.Fatalf("first Acquire failed: %v", err)
	}
	if !acquired {
		t.Fatal("expected first acquire to succeed")
	}

	acquired, err = manager.Acquire(ctx, 1, "key-1")
	if err != nil {
		t.Fatalf("second Acquire failed: %v", err)
	}
	if acquired {
		t.Fatal("expected second acquire to fail")
	}
}

func TestIdempotencyMarkSuccess(t *testing.T) {
	manager, _, client := newTestIdempotencyManager(t)
	ctx := context.Background()

	acquired, err := manager.Acquire(ctx, 1, "key-1")
	if err != nil {
		t.Fatalf("Acquire failed: %v", err)
	}
	if !acquired {
		t.Fatal("expected acquire to succeed")
	}

	if err := manager.MarkSuccess(ctx, 1, "key-1", "ORD123"); err != nil {
		t.Fatalf("MarkSuccess failed: %v", err)
	}

	val, err := client.Get(ctx, "order:idempotency:1:key-1").Result()
	if err != nil {
		t.Fatalf("failed to get redis value: %v", err)
	}
	if val != "SUCCESS:ORD123" {
		t.Fatalf("unexpected redis value: got %q want %q", val, "SUCCESS:ORD123")
	}
}

func TestIdempotencyRelease(t *testing.T) {
	manager, _, _ := newTestIdempotencyManager(t)
	ctx := context.Background()

	acquired, err := manager.Acquire(ctx, 1, "key-1")
	if err != nil {
		t.Fatalf("first Acquire failed: %v", err)
	}
	if !acquired {
		t.Fatal("expected first acquire to succeed")
	}

	if err := manager.Release(ctx, 1, "key-1"); err != nil {
		t.Fatalf("Release failed: %v", err)
	}

	acquired, err = manager.Acquire(ctx, 1, "key-1")
	if err != nil {
		t.Fatalf("second Acquire failed: %v", err)
	}
	if !acquired {
		t.Fatal("expected acquire after release to succeed")
	}
}

func TestIdempotencyDifferentUserIsolation(t *testing.T) {
	manager, _, _ := newTestIdempotencyManager(t)
	ctx := context.Background()

	acquired, err := manager.Acquire(ctx, 1, "same-key")
	if err != nil {
		t.Fatalf("Acquire for user 1 failed: %v", err)
	}
	if !acquired {
		t.Fatal("expected acquire for user 1 to succeed")
	}

	acquired, err = manager.Acquire(ctx, 2, "same-key")
	if err != nil {
		t.Fatalf("Acquire for user 2 failed: %v", err)
	}
	if !acquired {
		t.Fatal("expected acquire for user 2 to succeed")
	}
}
