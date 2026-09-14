package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/massodo1993/service-example/order/internal/api/order/v1"
	inventoryClient "github.com/massodo1993/service-example/order/internal/client/grpc/inventory"
	payemntCleint "github.com/massodo1993/service-example/order/internal/client/grpc/payment"
	"github.com/massodo1993/service-example/order/internal/config"
	"github.com/massodo1993/service-example/order/internal/migrator"
	orderRepository "github.com/massodo1993/service-example/order/internal/repository/order"
	orderService "github.com/massodo1993/service-example/order/internal/service/order"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

const (
	shutdownTimeout = 10 * time.Second
	handlerTimeout  = 30 * time.Second
)

const configPath = "./deploy/compose/order/.env"

func main() {
	if err := config.Load(configPath); err != nil {
		panic(fmt.Errorf("не удалось загрузить конфиг: %w", err))
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("не удалось подключиться к базе данных: %v\n", err)
		return
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("база данных недоступна: %v\n", err)
		return
	}

	migrationsDir := config.AppConfig().Postgres.MigrationsDir()
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("ошибка миграции базы данных: %v\n", err)
		return
	}

	inventoryConn, err := grpc.NewClient(
		config.AppConfig().InventoryGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("не удалось подключиться к inventory: %v\n", err)
		return
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("не удалось закрыть соединение с inventory: %v\n", cerr)
		}
	}()

	paymentConn, err := grpc.NewClient(
		config.AppConfig().PaymentGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("не удалось подключиться к payment: %v\n", err)
		return
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("не удалось закрыть соединение с payment: %v\n", cerr)
		}
	}()

	repo := orderRepository.NewRepository(pool)
	inventoryClient := inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(inventoryConn))
	paymentClient := payemntCleint.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))
	service := orderService.NewService(repo, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderAPI.NewApi(service))
	if err != nil {
		log.Printf("не удалось создать OpenAPI сервер: %v\n", err)
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(handlerTimeout))
	router.Mount("/", orderServer)

	server := &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           router,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	go func() {
		log.Printf("http order server listen on %s\n", config.AppConfig().OrderHTTP.Address())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("сервер остановлен с ошибкой: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("не удалось штатно остановить сервер: %v\n", err)
	}

	log.Println("server order stop")
}
