package state

import (
	"context"
	"fmt"
	"time"

	"tracking/internal/config"

	"github.com/redis/go-redis/v9"
)

// StateManager defines the interface for state management operations.
type StateManager interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisStore implements StateManager using Redis.
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new RedisStore using the provided configuration.
func NewRedisStore(cfg config.RedisConfig) (*RedisStore, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Println(addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisStore{client: client}, nil
}

// Set stores a value in Redis with an optional expiration.
func (s *RedisStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis.
func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

// Delete removes a value from Redis.
func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis.
func (s *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	n, err := s.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Close closes the Redis client.
func (s *RedisStore) Close() error {
	return s.client.Close()
}

var Redis *RedisStore

func Init(cfg config.RedisConfig) {
	var err error
	fmt.Printf("Starting State Module Empowered By Redis As Cache Layer: \t%+v\n", cfg.Host)

	Redis, err = NewRedisStore(cfg)
	if err != nil {
		panic(err)
	}
}
