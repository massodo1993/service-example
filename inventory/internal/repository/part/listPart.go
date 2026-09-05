package part

import (
	"context"

	"github.com/google/uuid"
	"github.com/massodo1993/service-example/inventory/internal/model"
	repoConverter "github.com/massodo1993/service-example/inventory/internal/repository/converter"
	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) ListParts(ctx context.Context, filters model.PartsFilter) ([]model.Part, error) {
	filter := filterToBson(filters)

	collect, err := r.mongo.Find(ctx, filter)
	if err != nil {
		return []model.Part{}, err
	}
	defer collect.Close(ctx)

	var parts []repoModel.Part
	if err = collect.All(ctx, &parts); err != nil {
		return []model.Part{}, err
	}

	return repoConverter.ToDomainParts(parts), nil
}

func filterToBson(filter model.PartsFilter) bson.M {
	m := bson.M{}

	if len(filter.UUIDs) > 0 {
		parsedUUIDs := make([]uuid.UUID, 0, len(filter.UUIDs))
		for _, u := range filter.UUIDs {
			parsed, err := uuid.Parse(u)
			if err != nil {
				continue
			}
			parsedUUIDs = append(parsedUUIDs, parsed)
		}
		m["part_uuid"] = bson.M{"$in": parsedUUIDs}
	}
	if len(filter.Names) > 0 {
		m["name"] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		m["category_type"] = bson.M{"$in": filter.Categories}
	}
	if len(filter.ManufacturerCountries) > 0 {
		m["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}
	if len(filter.Tags) > 0 {
		m["tags"] = bson.M{"$in": filter.Tags}
	}

	return m
}
