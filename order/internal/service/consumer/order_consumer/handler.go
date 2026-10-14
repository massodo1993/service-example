package order_consumer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
)

// ShipAssembledHandler переводит заказ в COMPLETED после сборки корабля.
func (s *service) ShipAssembledHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "не удалось декодировать ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "обработка ShipAssembled",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	orderUUID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		return fmt.Errorf("некорректный order_uuid %q: %w", event.OrderUUID, err)
	}

	if err := s.orderService.CompleteOrder(ctx, orderUUID); err != nil {
		logger.Error(ctx, "не удалось завершить заказ", zap.Error(err))
		return err
	}

	return nil
}
