package part

import (
	"context"
	"maps"
	"slices"

	"github.com/massodo1993/service-example/inventory/internal/model"
	repoConverter "github.com/massodo1993/service-example/inventory/internal/repository/converter"
	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filters model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoFilter := repoConverter.ToRepoFilter(filters)
	result := maps.Clone(r.data)
	if filters.IsEmpty() {
		parts := mapToSlice(result)
		return parts, nil
	}

	if len(filters.UUIDs) > 0 {
		maps.DeleteFunc(result, func(uuid string, _ repoModel.Part) bool {
			return !slices.Contains(repoFilter.UUIDs, uuid)
		})
	}
	if len(filters.Names) > 0 {
		maps.DeleteFunc(result, func(_ string, value repoModel.Part) bool {
			return !slices.Contains(repoFilter.Names, value.Name)
		})
	}
	if len(filters.Categories) > 0 {
		maps.DeleteFunc(result, func(_ string, value repoModel.Part) bool {
			return !slices.Contains(repoFilter.Categories, value.CategoryType)
		})
	}
	if len(filters.ManufacturerCountries) > 0 {
		maps.DeleteFunc(result, func(_ string, value repoModel.Part) bool {
			return !slices.Contains(repoFilter.ManufacturerCountries, value.Manufacturer.Country)
		})
	}
	if len(filters.Tags) > 0 {
		maps.DeleteFunc(result, func(_ string, value repoModel.Part) bool {
			for _, tag := range value.Tags {
				if slices.Contains(filters.Tags, tag) {
					return false
				}
			}
			return true
		})
	}

	parts := mapToSlice(result)
	return parts, nil
}

func mapToSlice(m map[string]repoModel.Part) []model.Part {
	result := make([]model.Part, 0, len(m))
	for _, v := range m {
		result = append(result, repoConverter.ToDomainPart(v))
	}
	return result
}
