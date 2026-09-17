package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/massodo1993/service-example/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

const (
	collectionName = "parts"
)

type repository struct {
	mongo *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		mongo: db.Collection(collectionName),
	}
}
