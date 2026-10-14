package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/massodo1993/service-example/notification/internal/converter/kafka"
	def "github.com/massodo1993/service-example/notification/internal/service"
	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	consumer        kafka.Consumer
	decoder         kafkaConverter.ShipAssembledDecoder
	telegramService def.TelegramService
}

func NewService(consumer kafka.Consumer, decoder kafkaConverter.ShipAssembledDecoder, telegramService def.TelegramService) *service {
	return &service{
		consumer:        consumer,
		decoder:         decoder,
		telegramService: telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "запуск консьюмера order assembled")

	if err := s.consumer.Consume(ctx, s.ShipAssembledHandler); err != nil {
		logger.Error(ctx, "ошибка чтения из топика", zap.String("service", "order_assembled_consumer"), zap.Error(err))
		return err
	}

	return nil
}
