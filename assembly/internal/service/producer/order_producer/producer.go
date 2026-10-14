package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/massodo1993/service-example/assembly/internal/model"
	def "github.com/massodo1993/service-example/assembly/internal/service"
	"github.com/massodo1993/service-example/platform/logger"
	"github.com/massodo1993/service-example/platform/pkg/kafka"
	eventsV1 "github.com/massodo1993/service-example/shared/pkg/proto/events/v1"
)

var _ def.OrderProducerService = (*service)(nil)

type service struct {
	orderAssembledProducer kafka.Producer
}

func NewService(orderAssembledProducer kafka.Producer) *service {
	return &service{orderAssembledProducer: orderAssembledProducer}
}

func (p *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	payload, err := proto.Marshal(&eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	})
	if err != nil {
		logger.Error(ctx, "не удалось сериализовать ShipAssembled", zap.Error(err))
		return err
	}

	if err := p.orderAssembledProducer.Send(ctx, []byte(event.OrderUUID), payload); err != nil {
		logger.Error(ctx, "не удалось опубликовать ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
