package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/massodo1993/service-example/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUuid uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	transactionUuid := uuid.New() //To-do делать реквест в пеймент
	err := s.orderRepository.PayOrder(ctx, orderUuid.String(), transactionUuid.String(), paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUuid, nil
}
