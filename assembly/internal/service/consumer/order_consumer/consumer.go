package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/massodo1993/service-example/assembly/internal/converter/kafka"
	def "github.com/massodo1993/service-example/assembly/internal/service"
	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	orderProducer     def.OrderProducerService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	orderProducer def.OrderProducerService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		orderProducer:     orderProducer,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "запуск консьюмера order paid")

	if err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler); err != nil {
		logger.Error(ctx, "ошибка чтения из топика order.paid", zap.Error(err))
		return err
	}

	return nil
}
