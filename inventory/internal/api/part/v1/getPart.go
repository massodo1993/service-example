package v1

import (
	"context"

	"github.com/massodo1993/service-example/inventory/internal/converter"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, request *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, request.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.ToProtoPart(part),
	}, nil
}
