package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/massodo1993/service-example/order/internal/api/order/v1"
	inventoryClient "github.com/massodo1993/service-example/order/internal/client/grpc/inventory"
	payemntCleint "github.com/massodo1993/service-example/order/internal/client/grpc/payment"
	orderRepository "github.com/massodo1993/service-example/order/internal/repository/order"
	orderService "github.com/massodo1993/service-example/order/internal/service/order"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50053"
	orderHost            = "localhost"
	orderPort            = 8080

	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	handlerTimeout    = 30 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddr,
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
		paymentServiceAddr,
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

	repo := orderRepository.NewRepository()
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
		Addr:              net.JoinHostPort(orderHost, strconv.Itoa(orderPort)),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("http order server listen on %d\n", orderPort)
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
