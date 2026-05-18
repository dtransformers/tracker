package state

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"tracking/internal/config"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisStore(t *testing.T) {
	// Start miniredis
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	// Parse host and port from miniredis address
	parts := strings.Split(mr.Addr(), ":")
	host := parts[0]
	port, _ := strconv.Atoi(parts[1])

	cfg := config.RedisConfig{
		Host:     host,
		Port:     port,
		Database: 0,
	}

	store, err := NewRedisStore(cfg)
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		err := store.Set(ctx, "key1", "value1", 0)
		assert.NoError(t, err)

		val, err := store.Get(ctx, "key1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)
	})

	t.Run("Get Non-Existent", func(t *testing.T) {
		_, err := store.Get(ctx, "non-existent")
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("Exists", func(t *testing.T) {
		err := store.Set(ctx, "key2", "value2", 0)
		require.NoError(t, err)

		exists, err := store.Exists(ctx, "key2")
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = store.Exists(ctx, "non-existent")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Set(ctx, "key3", "value3", 0)
		require.NoError(t, err)

		err = store.Delete(ctx, "key3")
		assert.NoError(t, err)

		exists, err := store.Exists(ctx, "key3")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Expiration", func(t *testing.T) {
		err := store.Set(ctx, "key4", "value4", 100*time.Millisecond)
		assert.NoError(t, err)

		val, err := store.Get(ctx, "key4")
		assert.NoError(t, err)
		assert.Equal(t, "value4", val)

		// Fast-forward time in miniredis
		mr.FastForward(200 * time.Millisecond)

		_, err = store.Get(ctx, "key4")
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})
}
