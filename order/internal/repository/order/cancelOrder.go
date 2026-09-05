package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/massodo1993/service-example/order/internal/model"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

func (r *repository) CancelOrder(ctx context.Context, orderUuid string) error {
	tag, err := r.connect.Exec(ctx,
		"UPDATE orders SET status = $1 WHERE order_uuid = $2 AND status = $3",
		repoModel.STATUS_CANCELLED, orderUuid, repoModel.STATUS_PENDING_PAYMENT,
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
		if err != nil {
			return err
		}
		return model.ErrOrderCancelConflict
	}

	return nil
}
