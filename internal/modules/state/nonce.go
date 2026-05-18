package state

import (
	"context"
	"fmt"
	"time"
)

type NonceStore struct {
	store *RedisStore
}

func NewNonceStore(store *RedisStore) *NonceStore {
	return &NonceStore{store: store}
}

func (s *NonceStore) CheckAndStore(publicKey, nonce string, expiry time.Time) bool {
	ctx := context.Background()
	key := fmt.Sprintf("nonce:%s:%s", publicKey, nonce)

	duration := time.Until(expiry)
	if duration <= 0 {
		return false
	}

	success, err := s.store.client.SetNX(ctx, key, "1", duration).Result()
	if err != nil {
		return false
	}
	return success
}
