package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUuid uuid.UUID) (model.User, error) {
	userUUID, err := s.sessionRepository.Get(ctx, sessionUuid)
	if err != nil {
		return model.User{}, err
	}

	return s.userRepository.GetByUUID(ctx, userUUID)
}
