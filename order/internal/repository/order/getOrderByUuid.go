package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/massodo1993/service-example/order/internal/model"
	"github.com/massodo1993/service-example/order/internal/repository/converter"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

func (r *repository) GetOrderByUuid(ctx context.Context, orderUuid string) (model.Order, error) {

	rows, err := r.connect.Query(ctx, "SELECT order_uuid, user_uuid, parts_uuids, total_price, transaction_uuid, payment_method, status FROM orders WHERE order_uuid = $1", orderUuid)
	if err != nil {
		return model.Order{}, err
	}

	repoOrder, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repoModel.Order])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Order{}, model.ErrOrderNotFound
	}
	if err != nil {
		return model.Order{}, err
	}

	return converter.ToModelOrder(repoOrder)
}
