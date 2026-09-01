package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/massodo1993/service-example/order/internal/model"
)

func newOrder() model.Order {
	return model.NewOrder(uuid.New(), []uuid.UUID{uuid.New(), uuid.New()}, model.STATUS_PENDING_PAYMENT, 123.45)
}

func TestCreateAndGetOrder(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()
	order := newOrder()

	require.NoError(t, repo.CreateOrder(ctx, order))

	got, err := repo.GetOrderByUuid(ctx, order.OrderUUID.String())
	require.NoError(t, err)
	require.Equal(t, order.OrderUUID, got.OrderUUID)
	require.Equal(t, order.UserUUID, got.UserUUID)
	require.Equal(t, order.TotalPrice, got.TotalPrice)
	require.Equal(t, model.STATUS_PENDING_PAYMENT, got.Status)
}

func TestGetOrderNotFound(t *testing.T) {
	repo := NewRepository()

	_, err := repo.GetOrderByUuid(context.Background(), uuid.NewString())
	require.ErrorIs(t, err, model.ErrOrderNotFound)
}

func TestPayOrder(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()
	order := newOrder()
	require.NoError(t, repo.CreateOrder(ctx, order))

	txUUID := uuid.NewString()
	require.NoError(t, repo.PayOrder(ctx, order.OrderUUID.String(), txUUID, model.PM_CARD))

	got, err := repo.GetOrderByUuid(ctx, order.OrderUUID.String())
	require.NoError(t, err)
	require.Equal(t, model.STATUS_PAID, got.Status)
	require.NotNil(t, got.TransactionUUID)

	// повторная оплата — конфликт
	require.ErrorIs(t, repo.PayOrder(ctx, order.OrderUUID.String(), txUUID, model.PM_CARD), model.ErrOrderPayConflict)
}

func TestPayOrderNotFound(t *testing.T) {
	err := NewRepository().PayOrder(context.Background(), uuid.NewString(), uuid.NewString(), model.PM_CARD)
	require.ErrorIs(t, err, model.ErrOrderNotFound)
}

func TestCancelOrder(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()
	order := newOrder()
	require.NoError(t, repo.CreateOrder(ctx, order))

	require.NoError(t, repo.CancelOrder(ctx, order.OrderUUID.String()))

	got, err := repo.GetOrderByUuid(ctx, order.OrderUUID.String())
	require.NoError(t, err)
	require.Equal(t, model.STATUS_CANCELLED, got.Status)

	// повторная отмена — конфликт
	require.ErrorIs(t, repo.CancelOrder(ctx, order.OrderUUID.String()), model.ErrOrderCancelConflict)
}

func TestCancelOrderNotFound(t *testing.T) {
	err := NewRepository().CancelOrder(context.Background(), uuid.NewString())
	require.ErrorIs(t, err, model.ErrOrderNotFound)
}
