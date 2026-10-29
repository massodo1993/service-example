package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, userUuid uuid.UUID) (model.User, error) {
	return s.userRepository.GetByUUID(ctx, userUuid)
}
