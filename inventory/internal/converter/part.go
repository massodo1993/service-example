package converter

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	svcModel "github.com/massodo1993/service-example/inventory/internal/model"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

func ToDomainPart(part *inventoryv1.Part) (svcModel.Part, error) {
	id, err := uuid.Parse(part.GetUuid())
	if err != nil {
		return svcModel.Part{}, err
	}

	return svcModel.Part{
		UUID:          id,
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: int(part.GetStockQuantity()),
		Category:      svcModel.CategoryType(part.GetCategory()),
		Dimensions:    toDomainDimensions(part.GetDimensions()),
		Manufacturer:  toDomainManufacturer(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      toDomainMetadata(part.GetMetadata()),
		CreatedAt:     part.GetCreatedAt().AsTime(),
		UpdatedAt:     protoTimeToDomain(part.GetUpdatedAt()),
	}, nil
}

func ToDomainParts(parts []*inventoryv1.Part) ([]svcModel.Part, error) {
	result := make([]svcModel.Part, 0, len(parts))
	for _, p := range parts {
		domainPart, err := ToDomainPart(p)
		if err != nil {
			return nil, err
		}
		result = append(result, domainPart)
	}
	return result, nil
}

func ToProtoPart(part svcModel.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: int64(part.StockQuantity),
		Category:      inventoryv1.CategoryType(part.Category), //nolint:gosec // enum с фиксированным набором значений, переполнение невозможно
		Dimensions:    toProtoDimensions(part.Dimensions),
		Manufacturer:  toProtoManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      toProtoMetadata(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     domainTimeToProto(part.UpdatedAt),
	}
}

func ToProtoParts(parts []svcModel.Part) []*inventoryv1.Part {
	result := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, ToProtoPart(p))
	}
	return result
}

func ToDomainFilter(filter *inventoryv1.PartsFilter) svcModel.PartsFilter {
	return svcModel.PartsFilter{
		UUIDs:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            toDomainCategories(filter.GetCategories()),
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

func toDomainDimensions(d *inventoryv1.Dimensions) *svcModel.Dimensions {
	if d == nil {
		return nil
	}
	return &svcModel.Dimensions{Length: d.GetLength(), Width: d.GetWidth(), Height: d.GetHeight(), Weight: d.GetWeight()}
}

func toProtoDimensions(d *svcModel.Dimensions) *inventoryv1.Dimensions {
	if d == nil {
		return nil
	}
	return &inventoryv1.Dimensions{Length: d.Length, Width: d.Width, Height: d.Height, Weight: d.Weight}
}

func toDomainManufacturer(m *inventoryv1.Manufacturer) *svcModel.Manufacturer {
	if m == nil {
		return nil
	}
	return &svcModel.Manufacturer{Name: m.GetName(), Country: m.GetCountry(), Website: m.GetWebsite()}
}

func toProtoManufacturer(m *svcModel.Manufacturer) *inventoryv1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryv1.Manufacturer{Name: m.Name, Country: m.Country, Website: m.Website}
}

func toDomainMetadata(src map[string]*inventoryv1.Value) map[string]svcModel.Value {
	if src == nil {
		return nil
	}
	dst := make(map[string]svcModel.Value, len(src))
	for k, v := range src {
		dst[k] = toDomainValue(v)
	}
	return dst
}

func toDomainValue(v *inventoryv1.Value) svcModel.Value {
	switch kind := v.GetKind().(type) {
	case *inventoryv1.Value_StringValue:
		return svcModel.Value{String_value: &kind.StringValue}
	case *inventoryv1.Value_Int64Value:
		return svcModel.Value{Int64_value: &kind.Int64Value}
	case *inventoryv1.Value_DoubleValue:
		return svcModel.Value{Double_value: &kind.DoubleValue}
	case *inventoryv1.Value_BoolValue:
		return svcModel.Value{Bool_value: &kind.BoolValue}
	default:
		return svcModel.Value{}
	}
}

func toProtoMetadata(src map[string]svcModel.Value) map[string]*inventoryv1.Value {
	if src == nil {
		return nil
	}
	dst := make(map[string]*inventoryv1.Value, len(src))
	for k, v := range src {
		dst[k] = toProtoValue(v)
	}
	return dst
}

func toProtoValue(v svcModel.Value) *inventoryv1.Value {
	switch {
	case v.String_value != nil:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_StringValue{StringValue: *v.String_value}}
	case v.Int64_value != nil:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_Int64Value{Int64Value: *v.Int64_value}}
	case v.Double_value != nil:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_DoubleValue{DoubleValue: *v.Double_value}}
	case v.Bool_value != nil:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_BoolValue{BoolValue: *v.Bool_value}}
	default:
		return &inventoryv1.Value{}
	}
}

func toDomainCategories(categories []inventoryv1.CategoryType) []svcModel.CategoryType {
	result := make([]svcModel.CategoryType, len(categories))
	for i, c := range categories {
		result[i] = svcModel.CategoryType(c)
	}
	return result
}

func protoTimeToDomain(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	return lo.ToPtr(ts.AsTime())
}

func domainTimeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(lo.FromPtr(t))
}
