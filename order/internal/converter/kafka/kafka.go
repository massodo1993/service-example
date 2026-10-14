package kafka

import "github.com/massodo1993/service-example/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
