package app

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryv1API "github.com/massodo1993/service-example/inventory/internal/api/part/v1"
	"github.com/massodo1993/service-example/inventory/internal/config"
	"github.com/massodo1993/service-example/inventory/internal/repository"
	partRepositroy "github.com/massodo1993/service-example/inventory/internal/repository/part"
	"github.com/massodo1993/service-example/inventory/internal/service"
	partService "github.com/massodo1993/service-example/inventory/internal/service/part"
	"github.com/massodo1993/service-example/platform/closer"
	grpcAuth "github.com/massodo1993/service-example/platform/pkg/middleware/grpc"
	authv1 "github.com/massodo1993/service-example/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryv1API inventoryv1.InventoryServiceServer
	partService    service.PartService
	partRepositroy repository.PartRepository
	mongoDBClient  *mongo.Client
	mongoDBHandle  *mongo.Database

	iamConn         *grpc.ClientConn
	iamClient       authv1.AuthServiceClient
	authInterceptor *grpcAuth.AuthInterceptor
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PartV1API(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryv1API == nil {
		d.inventoryv1API = inventoryv1API.NewAPI(d.InventoryService(ctx))
	}

	return d.inventoryv1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.partRepositroy == nil {
		repo := partRepositroy.NewRepository(d.MongoDBHandle(ctx))

		if err := repo.Seed(ctx); err != nil {
			panic(fmt.Sprintf("failed to seed parts: %s\n", err.Error()))
		}

		d.partRepositroy = repo
	}

	return d.partRepositroy
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(config.AppConfig().MongoConfig.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(connectCtx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().MongoConfig.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) AuthInterceptor(ctx context.Context) *grpcAuth.AuthInterceptor {
	if d.authInterceptor == nil {
		d.authInterceptor = grpcAuth.NewAuthInterceptor(d.IamClient(ctx))
	}

	return d.authInterceptor
}

func (d *diContainer) IamClient(ctx context.Context) authv1.AuthServiceClient {
	if d.iamClient == nil {
		d.iamClient = authv1.NewAuthServiceClient(d.IamConn(ctx))
	}

	return d.iamClient
}

func (d *diContainer) IamConn(_ context.Context) *grpc.ClientConn {
	if d.iamConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IamGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to iam: %s\n", err.Error()))
		}

		closer.AddNamed("iam grpc connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamConn = conn
	}

	return d.iamConn
}
