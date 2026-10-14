package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) CompleteOrder(ctx context.Context, orderUuid uuid.UUID) error {
	return s.orderRepository.CompleteOrder(ctx, orderUuid.String())
}
