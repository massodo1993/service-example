package order

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/massodo1993/service-example/order/internal/model"
	"github.com/massodo1993/service-example/platform/logger"
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

	// Заказ уже оплачен, поэтому ошибку публикации события только логируем.
	err = s.orderProducer.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUUID:       uuid.New(),
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   paymentMethod.String(),
		TransactionUUID: transactionUuid,
	})
	if err != nil {
		logger.Error(ctx, "не удалось опубликовать OrderPaid", zap.Error(err))
	}

	return transactionUuid, nil
}
