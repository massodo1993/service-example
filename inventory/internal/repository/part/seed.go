package part

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"

	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
)

// Seed наполняет коллекцию демо-данными, если она ещё пустая.
//
//nolint:gosec // G404: math/rand достаточно для генерации тестовых деталей, крипто-стойкость не нужна
func (r *repository) Seed(ctx context.Context) error {
	count, err := r.mongo.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	names := []string{
		"Двигатель Р7", "Крыло L-100", "Иллюминатор X1", "Топливный бак FT-9",
		"Двигатель М3", "Крыло R-200", "Иллюминатор X2", "Топливный бак FT-12",
		"Двигатель N9", "Крыло S-300",
	}

	allTags := []string{"металл", "пластик", "новое", "б/у", "premium", "эконом", "сертифицировано"}
	categories := []repoModel.CategoryType{
		repoModel.CATEGORY_TYPE_ENGINE,
		repoModel.CATEGORY_TYPE_FUEL,
		repoModel.CATEGORY_TYPE_PORTHOLE,
		repoModel.CATEGORY_TYPE_WING,
	}
	countries := []string{"Россия", "Германия", "США", "Китай", "Япония"}

	docs := make([]interface{}, 0, 10)

	for i := 0; i < 10 && len(names) > 0; i++ {
		nameIdx := rand.IntN(len(names))
		name := names[nameIdx]
		names = append(names[:nameIdx], names[nameIdx+1:]...)

		tagsCount := rand.IntN(3) + 1
		tags := make([]string, tagsCount)
		for j := range tags {
			tags[j] = allTags[rand.IntN(len(allTags))]
		}

		part := repoModel.Part{
			Uuid:          uuid.New(),
			Name:          name,
			Description:   "Сгенерировано автоматически",
			Price:         100 + rand.Float64()*9900,
			StockQuantity: rand.IntN(100) + 1,
			CategoryType:  categories[rand.IntN(len(categories))],
			Dimensions: &repoModel.Dimensions{
				Length: rand.Float64() * 200,
				Width:  rand.Float64() * 100,
				Height: rand.Float64() * 100,
				Weight: rand.Float64() * 500,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    fmt.Sprintf("Manufacturer %d", i+1),
				Country: countries[rand.IntN(len(countries))],
				Website: fmt.Sprintf("https://manufacturer%d.example.com", i+1),
			},
			Tags:      tags,
			Metadata:  make(map[string]repoModel.Value),
			CreatedAt: time.Now(),
			UpdatedAt: lo.ToPtr(time.Now()),
		}
		docs = append(docs, part)
	}

	_, err = r.mongo.InsertMany(ctx, docs)
	return err
}
