package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentAPI "github.com/massodo1993/service-example/payment/internal/api/payment/v1"
	paymentService "github.com/massodo1993/service-example/payment/internal/service/payment"
	paymentv1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

const paymentGRPCPort = 50053

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", paymentGRPCPort))
	if err != nil {
		log.Printf("не удалось занять порт: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("не удалось закрыть listener: %v\n", cerr)
		}
	}()

	server := grpc.NewServer()

	service := paymentService.NewService()
	api := paymentAPI.NewAPI(service)

	paymentv1.RegisterPaymentServiceServer(server, api)
	reflection.Register(server)

	go func() {
		log.Printf("grpc payment server listen on %d\n", paymentGRPCPort)
		if err := server.Serve(lis); err != nil {
			log.Printf("сервер остановлен с ошибкой: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.GracefulStop()
	log.Println("server payment stop")
}
