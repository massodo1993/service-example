//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"

	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
)

// InsertTestPart — вставляет тестовую деталь в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (uuid.UUID, error) {
	partUUID := uuid.New()
	now := time.Now()

	part := repoModel.Part{
		Uuid:          partUUID,
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Sentence(10),
		Price:         100 + gofakeit.Float64Range(0, 9900),
		StockQuantity: gofakeit.Number(1, 100),
		CategoryType:  repoModel.CATEGORY_TYPE_ENGINE,
		Dimensions: &repoModel.Dimensions{
			Length: gofakeit.Float64Range(1, 200),
			Width:  gofakeit.Float64Range(1, 100),
			Height: gofakeit.Float64Range(1, 100),
			Weight: gofakeit.Float64Range(1, 500),
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags:      []string{"тест"},
		Metadata:  make(map[string]repoModel.Value),
		CreatedAt: now,
	}

	databaseName := databaseNameFromEnv()

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, part)
	if err != nil {
		return uuid.Nil, err
	}

	return partUUID, nil
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := databaseNameFromEnv()

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	return err
}

func databaseNameFromEnv() string {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	return databaseName
}
