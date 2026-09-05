package part

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/massodo1993/service-example/inventory/internal/model"
	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
)

var testRepo *repository

func TestMain(m *testing.M) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://inventory-service-user:inventory-service-password@localhost:27017/inventory-service?authSource=admin"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect from mongo: %v\n", cerr)
		}
	}()

	collection := client.Database("inventory-service").Collection("parts")
	testRepo = &repository{mongo: collection}

	os.Exit(m.Run())
}

func truncate(t *testing.T) {
	t.Helper()

	_, err := testRepo.mongo.DeleteMany(context.Background(), bson.M{})
	require.NoError(t, err)
}

func seedParts(t *testing.T, n int) []repoModel.Part {
	t.Helper()

	parts := make([]repoModel.Part, 0, n)
	docs := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		part := repoModel.Part{
			Uuid:          uuid.New(),
			Name:          "Деталь",
			Description:   "тестовая деталь",
			Price:         100,
			StockQuantity: 1,
			CategoryType:  repoModel.CATEGORY_TYPE_ENGINE,
			CreatedAt:     time.Now(),
		}
		parts = append(parts, part)
		docs = append(docs, part)
	}

	_, err := testRepo.mongo.InsertMany(context.Background(), docs)
	require.NoError(t, err)

	return parts
}

func TestListPartsNoFilter(t *testing.T) {
	truncate(t)
	seedParts(t, 10)

	parts, err := testRepo.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(t, err)
	require.Len(t, parts, 10)
}

func TestGetPartSuccess(t *testing.T) {
	truncate(t)
	want := seedParts(t, 1)[0]

	got, err := testRepo.GetPart(context.Background(), want.Uuid.String())
	require.NoError(t, err)
	require.Equal(t, want.Uuid, got.UUID)
	require.Equal(t, want.Name, got.Name)
}

func TestGetPartNotFound(t *testing.T) {
	truncate(t)

	_, err := testRepo.GetPart(context.Background(), uuid.NewString())
	require.ErrorIs(t, err, model.ErrPartNotFound)
}

func TestListPartsFilterByUUID(t *testing.T) {
	truncate(t)
	all := seedParts(t, 3)
	target := all[0]

	parts, err := testRepo.ListParts(context.Background(), model.PartsFilter{UUIDs: []string{target.Uuid.String()}})
	require.NoError(t, err)
	require.Len(t, parts, 1)
	require.Equal(t, target.Uuid, parts[0].UUID)
}

func TestListPartsFilterByName(t *testing.T) {
	truncate(t)
	all := seedParts(t, 3)
	target := all[0]

	parts, err := testRepo.ListParts(context.Background(), model.PartsFilter{Names: []string{target.Name}})
	require.NoError(t, err)
	require.NotEmpty(t, parts)
	for _, p := range parts {
		require.Equal(t, target.Name, p.Name)
	}
}

func TestListPartsFilterNoMatch(t *testing.T) {
	truncate(t)
	seedParts(t, 3)

	parts, err := testRepo.ListParts(context.Background(), model.PartsFilter{Names: []string{"нет такой детали"}})
	require.NoError(t, err)
	require.Empty(t, parts)
}
