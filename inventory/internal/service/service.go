package service

import (
	"context"

	"github.com/massodo1993/service-example/inventory/internal/model"
)

type PartService interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filters model.PartsFilter) ([]model.Part, error)
}
