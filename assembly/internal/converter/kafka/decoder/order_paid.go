package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/massodo1993/service-example/assembly/internal/model"
	eventsV1 "github.com/massodo1993/service-example/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("не удалось десериализовать protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.GetEventUuid(),
		OrderUUID:       pb.GetOrderUuid(),
		UserUUID:        pb.GetUserUuid(),
		PaymentMethod:   pb.GetPaymentMethod(),
		TransactionUUID: pb.GetTransactionUuid(),
	}, nil
}
