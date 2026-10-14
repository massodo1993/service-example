package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/massodo1993/service-example/order/internal/converter/kafka"
	def "github.com/massodo1993/service-example/order/internal/service"
	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	shipAssembledDecoder   kafkaConverter.ShipAssembledDecoder
	orderService           def.OrderService
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder,
	orderService def.OrderService,
) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		shipAssembledDecoder:   shipAssembledDecoder,
		orderService:           orderService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "запуск консьюмера order assembled")

	if err := s.orderAssembledConsumer.Consume(ctx, s.ShipAssembledHandler); err != nil {
		logger.Error(ctx, "ошибка чтения из топика order.assembled", zap.Error(err))
		return err
	}

	return nil
}
