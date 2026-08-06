package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
	inventoryV1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

// ToModelPart конвертирует деталь из ответа inventory в доменную модель.
func ToModelPart(part *inventoryV1.Part) (model.Part, error) {
	if part == nil {
		return model.Part{}, fmt.Errorf("inventory вернул пустую деталь")
	}

	partUUID, err := uuid.Parse(part.GetUuid())
	if err != nil {
		return model.Part{}, fmt.Errorf("некорректный uuid детали %q: %w", part.GetUuid(), err)
	}

	metadata := make(map[string]model.Value, len(part.GetMetadata()))
	for key, value := range part.GetMetadata() {
		metadata[key] = toModelValue(value)
	}

	var updatedAt *time.Time
	if part.GetUpdatedAt() != nil {
		t := part.GetUpdatedAt().AsTime()
		updatedAt = &t
	}

	return model.Part{
		PartUUID:      partUUID,
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      toModelCategory(part.GetCategory()),
		Dimensions:    toModelDimensions(part.GetDimensions()),
		Manufacturer:  toModelManufacturer(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      metadata,
		CreatedAt:     part.GetCreatedAt().AsTime(),
		UpdatedAt:     updatedAt,
	}, nil
}

// ToModelParts конвертирует список деталей, прерываясь на первой некорректной.
func ToModelParts(parts []*inventoryV1.Part) ([]model.Part, error) {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		converted, err := ToModelPart(part)
		if err != nil {
			return nil, err
		}

		result = append(result, converted)
	}

	return result, nil
}

func toModelCategory(category inventoryV1.CategoryType) model.Category {
	switch category {
	case inventoryV1.CategoryType_CATEGORY_TYPE_ENGINE:
		return model.CATEGORY_ENGINE
	case inventoryV1.CategoryType_CATEGORY_TYPE_FUEL:
		return model.CATEGORY_FUEL
	case inventoryV1.CategoryType_CATEGORY_TYPE_PORTHOLE:
		return model.CATEGORY_PORTHOLE
	case inventoryV1.CategoryType_CATEGORY_TYPE_WING:
		return model.CATEGORY_WING
	default:
		return model.CATEGORY_UNKNOWN
	}
}

func toModelDimensions(dimensions *inventoryV1.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &model.Dimensions{
		Length: dimensions.GetLength(),
		Width:  dimensions.GetWidth(),
		Height: dimensions.GetHeight(),
		Weight: dimensions.GetWeight(),
	}
}

func toModelManufacturer(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    manufacturer.GetName(),
		Country: manufacturer.GetCountry(),
		Website: manufacturer.GetWebsite(),
	}
}

func toModelValue(value *inventoryV1.Value) model.Value {
	var result model.Value
	if value == nil {
		return result
	}

	switch kind := value.GetKind().(type) {
	case *inventoryV1.Value_StringValue:
		v := kind.StringValue
		result.StringValue = &v
	case *inventoryV1.Value_Int64Value:
		v := kind.Int64Value
		result.Int64Value = &v
	case *inventoryV1.Value_DoubleValue:
		v := kind.DoubleValue
		result.DoubleValue = &v
	case *inventoryV1.Value_BoolValue:
		v := kind.BoolValue
		result.BoolValue = &v
	}

	return result
}
