package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// addSessionToUserSet хранит все сессии пользователя в Redis set (например, чтобы завершить их разом).
func (r *repository) addSessionToUserSet(ctx context.Context, userUuid, sessionUuid uuid.UUID, ttl time.Duration) error {
	key := userSessionsKeyPrefix + userUuid.String()

	if err := r.cache.SAdd(ctx, key, sessionUuid.String()); err != nil {
		return err
	}

	return r.cache.Expire(ctx, key, ttl)
}
