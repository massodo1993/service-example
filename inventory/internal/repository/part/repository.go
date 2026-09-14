package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/massodo1993/service-example/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

const (
	databaseName   = "inventory-service"
	collectionName = "parts"
)

type repository struct {
	client *mongo.Client
	mongo  *mongo.Collection
}

func NewRepository(client *mongo.Client) *repository {
	return &repository{
		client: client,
		mongo:  client.Database(databaseName).Collection(collectionName),
	}
}
