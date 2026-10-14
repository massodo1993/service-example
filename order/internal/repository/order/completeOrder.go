package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/massodo1993/service-example/order/internal/model"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

// CompleteOrder переводит оплаченный заказ в COMPLETED. Повторный вызов — no-op (идемпотентность Kafka).
func (r *repository) CompleteOrder(ctx context.Context, orderUuid string) error {
	tag, err := r.connect.Exec(ctx,
		"UPDATE orders SET status = $1 WHERE order_uuid = $2 AND status = $3",
		repoModel.STATUS_COMPLETED, orderUuid, repoModel.STATUS_PAID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		var exists int
		err := r.connect.QueryRow(ctx, "SELECT 1 FROM orders WHERE order_uuid = $1", orderUuid).Scan(&exists)
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrOrderNotFound
		}

		return err
	}

	return nil
}
