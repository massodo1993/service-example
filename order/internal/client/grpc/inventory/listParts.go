package inventory

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/client/converter"
	"github.com/massodo1993/service-example/order/internal/model"
	inventoryV1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, partsUUIDs []uuid.UUID) ([]model.Part, error) {
	uuids := make([]string, 0, len(partsUUIDs))
	for _, partUUID := range partsUUIDs {
		uuids = append(uuids, partUUID.String())
	}

	response, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("запрос деталей в inventory: %w", err)
	}

	parts, err := converter.ToModelParts(response.GetParts())
	if err != nil {
		return nil, fmt.Errorf("конвертация ответа inventory: %w", err)
	}

	return parts, nil
}
