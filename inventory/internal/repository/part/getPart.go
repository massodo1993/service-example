package part

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/massodo1993/service-example/inventory/internal/model"
	repoConverter "github.com/massodo1993/service-example/inventory/internal/repository/converter"
	repoModel "github.com/massodo1993/service-example/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetPart(ctx context.Context, partUUID string) (model.Part, error) {
	parsedUUID, err := uuid.Parse(partUUID)
	if err != nil {
		return model.Part{}, model.ErrPartNotFound
	}

	var part repoModel.Part

	err = r.mongo.FindOne(ctx, bson.M{"part_uuid": parsedUUID}).Decode(&part)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Part{}, model.ErrPartNotFound
	}
	if err != nil {
		return model.Part{}, err
	}

	return repoConverter.ToDomainPart(part), nil
}
