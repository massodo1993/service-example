package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/massodo1993/service-example/order/internal/model"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

func (r *repository) PayOrder(ctx context.Context, orderUuid, transactionUuid string, paymentMethod model.PaymentMethod) error {
	pm := repoModel.PaymentMethod(paymentMethod)

	tag, err := r.connect.Exec(ctx,
		"UPDATE orders SET status = $1, transaction_uuid = $2, payment_method = $3 WHERE order_uuid = $4 AND status = $5",
		repoModel.STATUS_PAID, transactionUuid, pm, orderUuid, repoModel.STATUS_PENDING_PAYMENT,
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
		return model.ErrOrderPayConflict
	}

	return nil
}
