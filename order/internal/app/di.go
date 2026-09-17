package app

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/massodo1993/service-example/order/internal/api/order/v1"
	"github.com/massodo1993/service-example/order/internal/client"
	inventoryClient "github.com/massodo1993/service-example/order/internal/client/grpc/inventory"
	paymentClient "github.com/massodo1993/service-example/order/internal/client/grpc/payment"
	"github.com/massodo1993/service-example/order/internal/config"
	"github.com/massodo1993/service-example/order/internal/repository"
	orderRepository "github.com/massodo1993/service-example/order/internal/repository/order"
	"github.com/massodo1993/service-example/order/internal/service"
	orderService "github.com/massodo1993/service-example/order/internal/service/order"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/pkg/migrator/pg"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API      orderV1.Handler
	orderService    service.OrderService
	orderRepository repository.OrderRepository

	pgPool *pgxpool.Pool

	inventoryConn   *grpc.ClientConn
	inventoryClient client.InventoryClient

	paymentConn   *grpc.ClientConn
	paymentClient client.PaymentClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderAPI.NewApi(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		pool, err := pgxpool.New(connectCtx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to postgres: %s\n", err.Error()))
		}

		if err := pool.Ping(connectCtx); err != nil {
			panic(fmt.Sprintf("failed to ping postgres: %s\n", err.Error()))
		}

		migrationsDir := config.AppConfig().Postgres.MigrationsDir()
		migratorRunner := pg.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migrationsDir)

		if err := migratorRunner.Up(); err != nil {
			panic(fmt.Sprintf("failed to run migrations: %s\n", err.Error()))
		}

		closer.AddNamed("postgres pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool
}

func (d *diContainer) InventoryClient(ctx context.Context) client.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(d.InventoryConn(ctx)))
	}

	return d.inventoryClient
}

func (d *diContainer) InventoryConn(_ context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory: %s\n", err.Error()))
		}

		closer.AddNamed("inventory grpc connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryConn = conn
	}

	return d.inventoryConn
}

func (d *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentClient.NewClient(paymentV1.NewPaymentServiceClient(d.PaymentConn(ctx)))
	}

	return d.paymentClient
}

func (d *diContainer) PaymentConn(_ context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment: %s\n", err.Error()))
		}

		closer.AddNamed("payment grpc connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentConn = conn
	}

	return d.paymentConn
}
