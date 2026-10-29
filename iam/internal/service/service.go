package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
)

type UserService interface {
	Register(ctx context.Context, info model.UserRegistrationInfo) (uuid.UUID, error)
	GetUser(ctx context.Context, userUuid uuid.UUID) (model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
	Whoami(ctx context.Context, sessionUuid uuid.UUID) (model.User, error)
}
