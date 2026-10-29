package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := s.userRepository.GetByLogin(ctx, login)
	if errors.Is(err, model.ErrUserNotFound) {
		return uuid.Nil, model.ErrInvalidCredentials
	}
	if err != nil {
		return uuid.Nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, model.ErrInvalidCredentials
	}

	sessionUUID := uuid.New()
	if err := s.sessionRepository.Create(ctx, sessionUUID, user.UUID, s.sessionTTL); err != nil {
		return uuid.Nil, err
	}

	return sessionUUID, nil
}
