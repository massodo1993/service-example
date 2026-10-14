package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/massodo1993/service-example/notification/internal/model"
	eventsV1 "github.com/massodo1993/service-example/shared/pkg/proto/events/v1"
)

type shipAssembledDecoder struct{}

func NewShipAssembledDecoder() *shipAssembledDecoder {
	return &shipAssembledDecoder{}
}

func (d *shipAssembledDecoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("не удалось десериализовать protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
