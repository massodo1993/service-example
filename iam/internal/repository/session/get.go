package session

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
)

// Get возвращает UUID пользователя, которому принадлежит сессия.
func (r *repository) Get(ctx context.Context, sessionUuid uuid.UUID) (uuid.UUID, error) {
	value, err := r.cache.Get(ctx, sessionKeyPrefix+sessionUuid.String())
	if errors.Is(err, redigo.ErrNil) {
		return uuid.Nil, model.ErrSessionNotFound
	}
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.ParseBytes(value)
}
