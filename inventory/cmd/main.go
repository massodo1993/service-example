package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/joho/godotenv"
	inventoryv1API "github.com/massodo1993/service-example/inventory/internal/api/part/v1"
	inventroyRepository "github.com/massodo1993/service-example/inventory/internal/repository/part"
	inventoryService "github.com/massodo1993/service-example/inventory/internal/service/part"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
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

	ctx := context.Background()

	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		log.Printf("MONGO_URI is empty\n")
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to mongo: %v\n", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from mongo: %v\n", cerr)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("mongo is not reachable: %v\n", err)
		return
	}

	collection := client.Database("inventory-service").Collection("parts")

	server := grpc.NewServer()

	repo := inventroyRepository.NewRepository(collection)
	if err := repo.Seed(ctx); err != nil {
		log.Printf("failed to seed parts: %v\n", err)
		return
	}

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
