package client

import (
	"context"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, partsUUIDs []uuid.UUID) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUuid, userUuid uuid.UUID, pm model.PaymentMethod) (uuid.UUID, error)
}
