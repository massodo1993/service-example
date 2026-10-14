package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/massodo1993/service-example/order/internal/model"
	def "github.com/massodo1993/service-example/order/internal/service"
	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
	eventsV1 "github.com/massodo1993/service-example/shared/pkg/proto/events/v1"
)

var _ def.OrderProducerService = (*service)(nil)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{orderPaidProducer: orderPaidProducer}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	payload, err := proto.Marshal(&eventsV1.OrderPaid{
		EventUuid:       event.EventUUID.String(),
		OrderUuid:       event.OrderUUID.String(),
		UserUuid:        event.UserUUID.String(),
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID.String(),
	})
	if err != nil {
		logger.Error(ctx, "не удалось сериализовать OrderPaid", zap.Error(err))
		return err
	}

	if err := p.orderPaidProducer.Send(ctx, []byte(event.OrderUUID.String()), payload); err != nil {
		logger.Error(ctx, "не удалось опубликовать OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
