package part

import (
	"context"

	"github.com/massodo1993/service-example/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filters model.PartsFilter) ([]model.Part, error) {
	parts, err := s.PartRepository.ListParts(ctx, filters)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
