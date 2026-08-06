package order

import (
	"context"

	"github.com/massodo1993/service-example/order/internal/model"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

func (r *repository) CancelOrder(_ context.Context, orderUuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, has := r.data[orderUuid]
	if !has {
		return model.ErrOrderNotFound
	}
	if order.Status != repoModel.STATUS_PENDING_PAYMENT {
		return model.ErrOrderCancelConflict
	}
	order.Status = repoModel.STATUS_CANCELLED

	return nil
}
