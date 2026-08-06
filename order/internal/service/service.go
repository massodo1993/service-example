package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/massodo1993/service-example/order/internal/model"
)

type OrderService interface {
	GetOrderByUuid(ctx context.Context, orderUuid uuid.UUID) (model.Order, error)
	CreateOrder(ctx context.Context, userUuid uuid.UUID, partsUuid []uuid.UUID) (model.Order, error)
	PayOrder(ctx context.Context, orderUuid uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
	CancelOrder(ctx context.Context, orderUuid uuid.UUID) error
}
