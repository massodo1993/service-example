package order

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/massodo1993/service-example/order/internal/model"
)

var testRepo *repository

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://order-service-user:order-service-password@localhost:5435/order-service")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	testRepo = &repository{connect: pool}

	os.Exit(m.Run())
}

func truncate(t *testing.T) {
	t.Helper()

	_, err := testRepo.connect.Exec(context.Background(), "TRUNCATE orders")
	require.NoError(t, err)
}

func newOrder() model.Order {
	return model.NewOrder(uuid.New(), []uuid.UUID{uuid.New(), uuid.New()}, model.STATUS_PENDING_PAYMENT, 123.45)
}

func TestCreateAndGetOrder(t *testing.T) {
	truncate(t)

	ctx := context.Background()
	order := newOrder()

	require.NoError(t, testRepo.CreateOrder(ctx, order))

	got, err := testRepo.GetOrderByUuid(ctx, order.OrderUUID.String())
	require.NoError(t, err)
	require.Equal(t, order.OrderUUID, got.OrderUUID)
	require.Equal(t, order.UserUUID, got.UserUUID)
	require.Equal(t, order.TotalPrice, got.TotalPrice)
	require.Equal(t, model.STATUS_PENDING_PAYMENT, got.Status)
}

func TestGetOrderNotFound(t *testing.T) {
	truncate(t)

	_, err := testRepo.GetOrderByUuid(context.Background(), uuid.NewString())
	require.ErrorIs(t, err, model.ErrOrderNotFound)
}

func TestPayOrder(t *testing.T) {
	truncate(t)

	ctx := context.Background()
	order := newOrder()
	require.NoError(t, testRepo.CreateOrder(ctx, order))

	txUUID := uuid.NewString()
	require.NoError(t, testRepo.PayOrder(ctx, order.OrderUUID.String(), txUUID, model.PM_CARD))

	got, err := testRepo.GetOrderByUuid(ctx, order.OrderUUID.String())
	require.NoError(t, err)
	require.Equal(t, model.STATUS_PAID, got.Status)
	require.NotNil(t, got.TransactionUUID)

	// повторная оплата — конфликт
	require.ErrorIs(t, testRepo.PayOrder(ctx, order.OrderUUID.String(), txUUID, model.PM_CARD), model.ErrOrderPayConflict)
}

func TestPayOrderNotFound(t *testing.T) {
	truncate(t)

	err := testRepo.PayOrder(context.Background(), uuid.NewString(), uuid.NewString(), model.PM_CARD)
	require.ErrorIs(t, err, model.ErrOrderNotFound)
}

func TestCancelOrder(t *testing.T) {
	truncate(t)

	ctx := context.Background()
	order := newOrder()
	require.NoError(t, testRepo.CreateOrder(ctx, order))

	require.NoError(t, testRepo.CancelOrder(ctx, order.OrderUUID.String()))

	got, err := testRepo.GetOrderByUuid(ctx, order.OrderUUID.String())
	require.NoError(t, err)
	require.Equal(t, model.STATUS_CANCELLED, got.Status)

	// повторная отмена — конфликт
	require.ErrorIs(t, testRepo.CancelOrder(ctx, order.OrderUUID.String()), model.ErrOrderCancelConflict)
}

func TestCancelOrderNotFound(t *testing.T) {
	truncate(t)

	err := testRepo.CancelOrder(context.Background(), uuid.NewString())
	require.ErrorIs(t, err, model.ErrOrderNotFound)
}
