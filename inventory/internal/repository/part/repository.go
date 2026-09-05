package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/massodo1993/service-example/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mongo *mongo.Collection
}

func NewRepository(con *mongo.Collection) *repository {
	return &repository{
		mongo: con,
	}
}
