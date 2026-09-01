package part

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/massodo1993/service-example/inventory/internal/model"
)

func TestListPartsNoFilter(t *testing.T) {
	repo := NewRepository()

	parts, err := repo.ListParts(context.Background(), model.PartsFilter{})
	require.NoError(t, err)
	require.Len(t, parts, 10) // NewRepository сидит 10 деталей
}

func TestGetPartSuccess(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()

	all, err := repo.ListParts(ctx, model.PartsFilter{})
	require.NoError(t, err)
	want := all[0]

	got, err := repo.GetPart(ctx, want.UUID.String())
	require.NoError(t, err)
	require.Equal(t, want.UUID, got.UUID)
	require.Equal(t, want.Name, got.Name)
}

func TestGetPartNotFound(t *testing.T) {
	_, err := NewRepository().GetPart(context.Background(), uuid.NewString())
	require.ErrorIs(t, err, model.ErrPartNotFound)
}

func TestListPartsFilterByUUID(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()

	all, err := repo.ListParts(ctx, model.PartsFilter{})
	require.NoError(t, err)
	target := all[0]

	parts, err := repo.ListParts(ctx, model.PartsFilter{UUIDs: []string{target.UUID.String()}})
	require.NoError(t, err)
	require.Len(t, parts, 1)
	require.Equal(t, target.UUID, parts[0].UUID)
}

func TestListPartsFilterByName(t *testing.T) {
	ctx := context.Background()
	repo := NewRepository()

	all, err := repo.ListParts(ctx, model.PartsFilter{})
	require.NoError(t, err)
	target := all[0]

	parts, err := repo.ListParts(ctx, model.PartsFilter{Names: []string{target.Name}})
	require.NoError(t, err)
	require.NotEmpty(t, parts)
	for _, p := range parts {
		require.Equal(t, target.Name, p.Name)
	}
}

func TestListPartsFilterNoMatch(t *testing.T) {
	parts, err := NewRepository().ListParts(context.Background(), model.PartsFilter{Names: []string{"нет такой детали"}})
	require.NoError(t, err)
	require.Empty(t, parts)
}
