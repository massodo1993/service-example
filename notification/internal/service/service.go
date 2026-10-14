package service

import (
	"context"

	"github.com/massodo1993/service-example/notification/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error
	SendShipAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error
}
