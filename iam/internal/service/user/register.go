package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *service) Register(ctx context.Context, info model.UserRegistrationInfo) (uuid.UUID, error) {
	if info.Info.Login == "" || info.Info.Email == "" || info.Password == "" {
		return uuid.Nil, model.ErrInvalidArgument
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("хэширование пароля: %w", err)
	}

	user, err := s.userRepository.Create(ctx, info.Info, string(hash))
	if err != nil {
		return uuid.Nil, err
	}

	return user.UUID, nil
}
