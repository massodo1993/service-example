package main

import (
	"context"
	"fmt"
	"log"
	"maps"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/google/uuid"
	inventoryv1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const inventoryGRPCPort = 50051

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	mu    sync.Mutex
	parts map[string]*inventoryv1.Part
}

func (is *inventoryService) GetPart(_ context.Context, request *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	is.mu.Lock()
	defer is.mu.Unlock()

	part, empty := is.parts[request.GetUuid()]
	if !empty {
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

func mockParts(service *inventoryService) {
	names := []string{
		"Двигатель Р7", "Крыло L-100", "Иллюминатор X1", "Топливный бак FT-9",
		"Двигатель М3", "Крыло R-200", "Иллюминатор X2", "Топливный бак FT-12",
		"Двигатель N9", "Крыло S-300",
	}

	allTags := []string{"металл", "пластик", "новое", "б/у", "premium", "эконом", "сертифицировано"}
	categories := []inventoryv1.CategoryType{
		inventoryv1.CategoryType_CATEGORY_TYPE_ENGINE,
		inventoryv1.CategoryType_CATEGORY_TYPE_FUEL,
		inventoryv1.CategoryType_CATEGORY_TYPE_PORTHOLE,
		inventoryv1.CategoryType_CATEGORY_TYPE_WING,
	}
	countries := []string{"Россия", "Германия", "США", "Китай", "Япония"}

	service.mu.Lock()
	defer service.mu.Unlock()

	for i := 0; i < 10 && len(names) > 0; i++ {
		nameIdx := rand.IntN(len(names))
		name := names[nameIdx]
		names = append(names[:nameIdx], names[nameIdx+1:]...)

		tagsCount := rand.IntN(3) + 1
		tags := make([]string, tagsCount)
		for j := range tags {
			tags[j] = allTags[rand.IntN(len(allTags))]
		}

		part := inventoryv1.Part{
			Uuid:          uuid.New().String(),
			Name:          name,
			Description:   "Сгенерировано автоматически",
			Price:         100 + rand.Float64()*9900,
			StockQuantity: int64(rand.IntN(100) + 1),
			Category:      categories[rand.IntN(len(categories))],
			Dimensions: &inventoryv1.Dimensions{
				Length: rand.Float64() * 200,
				Width:  rand.Float64() * 100,
				Height: rand.Float64() * 100,
				Weight: rand.Float64() * 500,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    fmt.Sprintf("Manufacturer %d", i+1),
				Country: countries[rand.IntN(len(countries))],
				Website: fmt.Sprintf("https://manufacturer%d.example.com", i+1),
			},
			Tags:      tags,
			Metadata:  make(map[string]*inventoryv1.Value),
			CreatedAt: timestamppb.Now(),
			UpdatedAt: timestamppb.Now(),
		}
		service.parts[part.Uuid] = &part
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", inventoryGRPCPort))
	if err != nil {
		log.Printf("faill listen: %v\n", err)
		return
	}

	server := grpc.NewServer()

	service := &inventoryService{
		parts: make(map[string]*inventoryv1.Part),
	}
	mockParts(service)

	inventoryv1.RegisterInventoryServiceServer(server, service)
	reflection.Register(server)

	go func() {
		log.Printf("grpc inventory server listen on %d\n", inventoryGRPCPort)
		err = server.Serve(lis)
		if err != nil {
			log.Printf("filed to server: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.GracefulStop()
	log.Println("server inventory stop")
}
