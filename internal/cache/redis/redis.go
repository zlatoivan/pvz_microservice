package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis struct in Redis
type Redis struct {
	mu         sync.RWMutex
	client     *redis.Client
	expiration time.Duration
}

// New initializes a new Redis
func New(addr string, expiration time.Duration) *Redis {
	return &Redis{
		mu: sync.RWMutex{},
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		expiration: expiration,
	}
}

// Close closes client
func (r *Redis) Close() error {
	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("r.client.Close: %w", err)
	}
	return nil
}

// Set sets a value in Redis by key
func (r *Redis) Set(ctx context.Context, key string, value []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.client.Set(ctx, key, value, r.expiration).Err()
	if err != nil {
		return fmt.Errorf("r.client.Set: %w", err)
	}

	return nil
}

// Get gets a value from Redis by key
func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("r.client.Get: %w", err)
	}

	return val, nil
}

// Delete deletes a value from Redis by key
func (r *Redis) Delete(ctx context.Context, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("r.client.Del: %w", err)
	}

	return nil
}
