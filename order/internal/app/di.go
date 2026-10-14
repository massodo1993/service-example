package app

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/massodo1993/service-example/order/internal/api/order/v1"
	"github.com/massodo1993/service-example/order/internal/client"
	inventoryClient "github.com/massodo1993/service-example/order/internal/client/grpc/inventory"
	paymentClient "github.com/massodo1993/service-example/order/internal/client/grpc/payment"
	"github.com/massodo1993/service-example/order/internal/config"
	kafkaConverter "github.com/massodo1993/service-example/order/internal/converter/kafka"
	"github.com/massodo1993/service-example/order/internal/converter/kafka/decoder"
	"github.com/massodo1993/service-example/order/internal/repository"
	orderRepository "github.com/massodo1993/service-example/order/internal/repository/order"
	"github.com/massodo1993/service-example/order/internal/service"
	orderConsumer "github.com/massodo1993/service-example/order/internal/service/consumer/order_consumer"
	orderService "github.com/massodo1993/service-example/order/internal/service/order"
	orderProducer "github.com/massodo1993/service-example/order/internal/service/producer/order_producer"
	"github.com/massodo1993/service-example/platform/closer"
	"github.com/massodo1993/service-example/platform/logger"
	wrappedKafka "github.com/massodo1993/service-example/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/massodo1993/service-example/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/massodo1993/service-example/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/massodo1993/service-example/platform/pkg/middleware/kafka"
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

	orderProducerService service.OrderProducerService
	consumerService      service.ConsumerService

	consumerGroup          sarama.ConsumerGroup
	orderAssembledConsumer wrappedKafka.Consumer
	shipAssembledDecoder   kafkaConverter.ShipAssembledDecoder

	syncProducer      sarama.SyncProducer
	orderPaidProducer wrappedKafka.Producer
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
			d.OrderProducerService(),
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

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderPaidProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = orderConsumer.NewService(
			d.OrderAssembledConsumer(),
			d.ShipAssembledDecoder(),
			d.OrderService(ctx),
		)
	}

	return d.consumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().OrderAssembledConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.shipAssembledDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderPaidProducer() wrappedKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer
}
