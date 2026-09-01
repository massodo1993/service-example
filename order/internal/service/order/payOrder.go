package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUuid uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.GetOrderByUuid(ctx, orderUuid)
	if err != nil {
		return uuid.Nil, err
	}

	transactionUuid, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	err = s.orderRepository.PayOrder(ctx, orderUuid.String(), transactionUuid.String(), paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUuid, nil
}
