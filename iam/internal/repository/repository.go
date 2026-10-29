package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info model.UserInfo, passwordHash string) (model.User, error)
	GetByUUID(ctx context.Context, userUuid uuid.UUID) (model.User, error)
	GetByLogin(ctx context.Context, login string) (model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, sessionUuid, userUuid uuid.UUID, ttl time.Duration) error
	Get(ctx context.Context, sessionUuid uuid.UUID) (uuid.UUID, error)
}
