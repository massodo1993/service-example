package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userUuid uuid.UUID, partsUuid []uuid.UUID) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, partsUuid)
	if err != nil {
		return model.Order{}, err
	}

	partsByUUID := make(map[uuid.UUID]model.Part, len(parts))
	for _, part := range parts {
		partsByUUID[part.PartUUID] = part
	}

	var totalPrice float64
	for _, partUUID := range partsUuid {
		part, found := partsByUUID[partUUID]
		if !found {
			return model.Order{}, fmt.Errorf("%w: %s", model.ErrPartNotFound, partUUID)
		}

		totalPrice += part.Price
	}

	order := model.NewOrder(userUuid, partsUuid, model.STATUS_PENDING_PAYMENT, totalPrice)
	err = s.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
