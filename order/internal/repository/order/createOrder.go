package order

import (
	"context"

	"github.com/massodo1993/service-example/order/internal/model"
	"github.com/massodo1993/service-example/order/internal/repository/converter"
)

func (r *repository) CreateOrder(ctx context.Context, order model.Order) error {
	repoOrder := converter.ToRepoOrder(order)
	_, err := r.connect.Exec(ctx,
		"INSERT INTO orders (order_uuid, user_uuid, parts_uuids, total_price, status) VALUES ($1, $2, $3, $4, $5)",
		repoOrder.OrderUUID, repoOrder.UserUUID, repoOrder.PartsUUIDs, repoOrder.TotalPrice, repoOrder.Status)

	if err != nil {
		return err
	}

	return nil
}
