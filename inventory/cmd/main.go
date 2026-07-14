package main

import (
	"context"
	"maps"
	"slices"
	"sync"

	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	mu    sync.Mutex
	parts map[string]*inventoryv1.Part
}

func (is *inventoryService) GetPart(_ context.Context, request *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	is.mu.Lock()
	defer is.mu.Unlock()

	part, empty := is.parts[request.GetUuid()]
	if empty {
		return nil, status.Errorf(codes.NotFound, "Part with UUID %s not found", request.GetUuid())
	}

	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

func (is *inventoryService) ListParts(_ context.Context, request *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	is.mu.Lock()
	defer is.mu.Unlock()
	result := maps.Clone(is.parts)
	if isEmpty(request.Filter) {
		parts := mapToSlice(result)
		return &inventoryv1.ListPartsResponse{Parts: parts}, nil
	}

	if len(request.Filter.GetUuids()) > 0 {
		maps.DeleteFunc(result, func(uuid string, _ *inventoryv1.Part) bool {
			return !slices.Contains(request.Filter.GetUuids(), uuid)
		})
	}
	if len(request.Filter.GetNames()) > 0 {
		maps.DeleteFunc(result, func(_ string, value *inventoryv1.Part) bool {
			return !slices.Contains(request.Filter.GetNames(), value.Name)
		})
	}
	if len(request.Filter.GetCategories()) > 0 {
		maps.DeleteFunc(result, func(_ string, value *inventoryv1.Part) bool {
			return !slices.Contains(request.Filter.GetCategories(), value.Category)
		})
	}
	if len(request.Filter.GetManufacturerCountries()) > 0 {
		maps.DeleteFunc(result, func(_ string, value *inventoryv1.Part) bool {
			return !slices.Contains(request.Filter.GetManufacturerCountries(), value.Manufacturer.String())
		})
	}
	if len(request.Filter.GetTags()) > 0 {
		maps.DeleteFunc(result, func(_ string, value *inventoryv1.Part) bool {
			for _, tag := range value.GetTags() {
				if slices.Contains(request.Filter.GetTags(), tag) {
					return false
				}
			}
			return true
		})
	}

	parts := mapToSlice(result)
	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func mapToSlice(m map[string]*inventoryv1.Part) []*inventoryv1.Part {
	result := make([]*inventoryv1.Part, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func isEmpty(f *inventoryv1.PartsFilter) bool {
	if f == nil {
		return true
	}
	return len(f.GetUuids()) == 0 &&
		len(f.GetNames()) == 0 &&
		len(f.GetCategories()) == 0 &&
		len(f.GetManufacturerCountries()) == 0 &&
		len(f.GetTags()) == 0
}

func main() {

}
