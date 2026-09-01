package repository

import (
	"context"

	"github.com/massodo1993/service-example/order/internal/model"
)

type OrderRepository interface {
	GetOrderByUuid(ctx context.Context, orderUuid string) (model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
	PayOrder(ctx context.Context, orderUuid, transactionUuid string, paymentMethod model.PaymentMethod) error
	CancelOrder(ctx context.Context, orderUuid string) error
}
