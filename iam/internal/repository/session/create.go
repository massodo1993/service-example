package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Create сохраняет сессию в Redis с TTL и добавляет её в множество сессий пользователя.
func (r *repository) Create(ctx context.Context, sessionUuid, userUuid uuid.UUID, ttl time.Duration) error {
	if err := r.cache.SetWithTTL(ctx, sessionKeyPrefix+sessionUuid.String(), userUuid.String(), ttl); err != nil {
		return err
	}

	return r.addSessionToUserSet(ctx, userUuid, sessionUuid, ttl)
}
