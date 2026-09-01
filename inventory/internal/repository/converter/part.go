package converter

import (
	svcModel "github.com/massodo1993/service-example/inventory/internal/model"
	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
)

// Из домена (internal/model) → в модель хранилища (internal/repository/model).
// Используется сервисным слоем перед вызовом repository.
func ToRepoPart(part svcModel.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		CategoryType:  repoModel.CategoryType(part.Category),
		Dimensions:    toRepoDimensions(part.Dimensions),
		Manufacturer:  toRepoManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      toRepoMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
		DeletedAt:     part.DeletedAt,
	}
}

func ToRepoParts(parts []svcModel.Part) []repoModel.Part {
	repoParts := make([]repoModel.Part, 0, len(parts))
	for _, value := range parts {
		repoParts = append(repoParts, ToRepoPart(value))
	}
	return repoParts
}

func ToRepoFilter(filter svcModel.PartsFilter) repoModel.PartsFilter {
	return repoModel.PartsFilter{
		UUIDs:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            toRepoCategories(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// Из модели хранилища (internal/repository/model) → обратно в домен (internal/model).
// Используется сервисным слоем после вызова repository.
func ToDomainPart(part repoModel.Part) svcModel.Part {
	return svcModel.Part{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      svcModel.CategoryType(part.CategoryType),
		Dimensions:    toDomainDimensions(part.Dimensions),
		Manufacturer:  toDomainManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      toDomainMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
		DeletedAt:     part.DeletedAt,
	}
}

func ToDomainParts(parts []repoModel.Part) []svcModel.Part {
	domainParts := make([]svcModel.Part, 0, len(parts))
	for _, value := range parts {
		domainParts = append(domainParts, ToDomainPart(value))
	}
	return domainParts
}

func toRepoMetadata(src map[string]svcModel.Value) map[string]repoModel.Value {
	if src == nil {
		return nil
	}
	dst := make(map[string]repoModel.Value, len(src))
	for k, v := range src {
		dst[k] = repoModel.Value(v)
	}
	return dst
}

func toDomainMetadata(src map[string]repoModel.Value) map[string]svcModel.Value {
	if src == nil {
		return nil
	}
	dst := make(map[string]svcModel.Value, len(src))
	for k, v := range src {
		dst[k] = svcModel.Value(v)
	}
	return dst
}

func toRepoCategories(categories []svcModel.CategoryType) []repoModel.CategoryType {
	result := make([]repoModel.CategoryType, len(categories))
	for i, c := range categories {
		result[i] = repoModel.CategoryType(c)
	}
	return result
}

func toRepoDimensions(d *svcModel.Dimensions) *repoModel.Dimensions {
	if d == nil {
		return nil
	}

	return &repoModel.Dimensions{Length: d.Length, Width: d.Width, Height: d.Height, Weight: d.Weight}
}

func toDomainDimensions(d *repoModel.Dimensions) *svcModel.Dimensions {
	if d == nil {
		return nil
	}

	return &svcModel.Dimensions{Length: d.Length, Width: d.Width, Height: d.Height, Weight: d.Weight}
}

func toRepoManufacturer(m *svcModel.Manufacturer) *repoModel.Manufacturer {
	if m == nil {
		return nil
	}

	return &repoModel.Manufacturer{Name: m.Name, Country: m.Country, Website: m.Website}
}

func toDomainManufacturer(m *repoModel.Manufacturer) *svcModel.Manufacturer {
	if m == nil {
		return nil
	}

	return &svcModel.Manufacturer{Name: m.Name, Country: m.Country, Website: m.Website}
}
