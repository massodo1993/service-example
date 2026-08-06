package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/massodo1993/service-example/order/internal/model"
)

func (s *service) GetOrderByUuid(ctx context.Context, orderUuid uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.GetOrderByUuid(ctx, orderUuid.String())

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
