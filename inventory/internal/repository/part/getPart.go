package part

import (
	"context"

	"github.com/massodo1993/service-example/inventory/internal/model"
	repoConverter "github.com/massodo1993/service-example/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.data[uuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.ToDomainPart(repoPart), nil
}
