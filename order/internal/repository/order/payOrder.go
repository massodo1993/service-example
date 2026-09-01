package order

import (
	"context"

	"github.com/massodo1993/service-example/order/internal/model"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

func (r *repository) PayOrder(_ context.Context, orderUuid, transactionUuid string, paymentMethod model.PaymentMethod) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, has := r.data[orderUuid]
	if !has {
		return model.ErrOrderNotFound
	}

	if order.Status != repoModel.STATUS_PENDING_PAYMENT {
		return model.ErrOrderPayConflict
	}

	pm := repoModel.PaymentMethod(paymentMethod)
	order.TransactionUUID = &transactionUuid
	order.PaymentMethod = &pm
	order.Status = repoModel.STATUS_PAID

	return nil
}
