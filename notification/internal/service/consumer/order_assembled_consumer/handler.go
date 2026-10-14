package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
)

func (s *service) ShipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.decoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "не удалось декодировать ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "обработка ShipAssembled",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
	)

	return s.telegramService.SendShipAssembledNotification(ctx, event)
}
