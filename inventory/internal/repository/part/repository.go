package part

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	def "github.com/massodo1993/service-example/inventory/internal/repository"
	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	r := &repository{
		data: make(map[string]repoModel.Part),
	}
	r.mockParts()
	return r
}

// mockParts наполняет in-memory хранилище демо-данными.
//
//nolint:gosec // G404: math/rand достаточно для генерации тестовых деталей, крипто-стойкость не нужна
func (r *repository) mockParts() {
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

	r.mu.Lock()
	defer r.mu.Unlock()

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
		r.data[part.Uuid.String()] = part
	}
}
