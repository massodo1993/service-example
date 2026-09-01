package order

import (
	"context"

	"github.com/massodo1993/service-example/order/internal/model"
	repoConverter "github.com/massodo1993/service-example/order/internal/repository/converter"
)

func (r *repository) GetOrderByUuid(_ context.Context, orderUuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, has := r.data[orderUuid]
	if !has {
		return model.Order{}, model.ErrOrderNotFound
	}
	modelOrder, err := repoConverter.ToModelOrder(*order)
	if err != nil {
		return model.Order{}, err
	}

	return modelOrder, nil
}
