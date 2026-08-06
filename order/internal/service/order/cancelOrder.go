package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) CancelOrder(ctx context.Context, orderUuid uuid.UUID) error {
	return s.orderRepository.CancelOrder(ctx, orderUuid.String())
}
