package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	paymentv1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const payemntGRPCPort = 50053

type payemntService struct {
	paymentv1.UnimplementedPayemntServiceServer

	mu       sync.Mutex
	payments map[string]*paymentv1.PaymentMethod
}

func (ps *payemntService) PayOrder(_ context.Context, request *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	uuid := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", uuid)
	return &paymentv1.PayOrderResponse{TransactionUuid: uuid.String()}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", payemntGRPCPort))
	if err != nil {
		log.Printf("faill listen: %v\n", err)
		return
	}

	server := grpc.NewServer()

	service := &payemntService{
		payments: make(map[string]*paymentv1.PaymentMethod),
	}

	paymentv1.RegisterPayemntServiceServer(server, service)
	reflection.Register(server)

	go func() {
		log.Printf("grpc payment server listen on %d\n", payemntGRPCPort)
		err = server.Serve(lis)
		if err != nil {
			log.Printf("filed to server: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.GracefulStop()
	log.Println("server stop")
}
