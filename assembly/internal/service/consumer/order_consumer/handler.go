package order_consumer

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/massodo1993/service-example/assembly/internal/model"
	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
)

// maxBuildSec — верхняя граница имитации сборки (1..maxBuildSec секунд).
const maxBuildSec = 10

// OrderPaidHandler имитирует сборку корабля и публикует ShipAssembled.
// Ожидание выполняется прямо в обработчике: offset коммитится только после публикации события.
func (s *service) OrderPaidHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "не удалось декодировать OrderPaid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "обработка OrderPaid",
		zap.String("topic", msg.Topic),
		zap.Int32("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	buildSec := int64(rand.IntN(maxBuildSec) + 1)

	select {
	case <-time.After(time.Duration(buildSec) * time.Second):
	case <-ctx.Done():
		return ctx.Err()
	}

	return s.orderProducer.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildSec,
	})
}
