package order

import (
	"context"

	"github.com/massodo1993/service-example/order/internal/model"
	"github.com/massodo1993/service-example/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	repoOrder := converter.ToRepoOrder(order)
	r.data[order.OrderUUID.String()] = &repoOrder

	return nil
}
