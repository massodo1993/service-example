package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryv1API "github.com/massodo1993/service-example/inventory/internal/api/part/v1"
	"github.com/massodo1993/service-example/inventory/internal/config"
	inventroyRepository "github.com/massodo1993/service-example/inventory/internal/repository/part"
	inventoryService "github.com/massodo1993/service-example/inventory/internal/service/part"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("не удалось загрузить конфиг: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().MongoConfig.URI()))
	if err != nil {
		log.Printf("не получилось подключиться к mongodb: %v\n", err)
		return
	}
	defer func() {
		if cerr := mongoClient.Disconnect(context.Background()); cerr != nil {
			log.Printf("не получилось отключиться от mongodb: %v\n", cerr)
		}
	}()

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Printf("mongodb не отвечает на ping: %v\n", err)
		return
	}
	log.Println("подключено к mongodb")

	lis, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
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

	repo := inventroyRepository.NewRepository(mongoClient)
	if err := repo.Seed(ctx); err != nil {
		log.Printf("failed to seed parts: %v\n", err)
		return
	}

	service := inventoryService.NewService(repo)
	api := inventoryv1API.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(server, api)

	reflection.Register(server)

	go func() {
		log.Printf("grpc inventory server listen on %s\n", config.AppConfig().InventoryGRPC.Address())
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
