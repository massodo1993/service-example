package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryv1API "github.com/massodo1993/service-example/inventory/internal/api/part/v1"
	inventroyRepository "github.com/massodo1993/service-example/inventory/internal/repository/part"
	inventoryService "github.com/massodo1993/service-example/inventory/internal/service/part"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const inventoryGRPCPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", inventoryGRPCPort))
	if err != nil {
		log.Printf("faill listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener %v\n", cerr)
		}
	}()

	server := grpc.NewServer()

	repo := inventroyRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryv1API.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(server, api)

	reflection.Register(server)

	go func() {
		log.Printf("grpc inventory server listen on %d\n", inventoryGRPCPort)
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
	log.Println("server inventory stop")
}
