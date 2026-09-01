package v1

import (
	"context"

	"github.com/massodo1993/service-example/inventory/internal/converter"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, request *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.partService.ListParts(ctx, converter.ToDomainFilter(request.Filter))
	if err != nil {
		return nil, err
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.ToProtoParts(parts),
	}, nil
}
