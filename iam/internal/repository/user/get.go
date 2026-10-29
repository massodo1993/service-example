package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/massodo1993/service-example/iam/internal/model"
	"github.com/massodo1993/service-example/iam/internal/repository/converter"
	repoModel "github.com/massodo1993/service-example/iam/internal/repository/model"
)

const selectUser = "SELECT uuid, login, email, password_hash, notification_methods, created_at, updated_at FROM users WHERE "

func (r *repository) GetByUUID(ctx context.Context, userUuid uuid.UUID) (model.User, error) {
	return r.get(ctx, selectUser+"uuid = $1", userUuid)
}

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	return r.get(ctx, selectUser+"login = $1", login)
}

func (r *repository) get(ctx context.Context, query string, arg any) (model.User, error) {
	rows, err := r.connect.Query(ctx, query, arg)
	if err != nil {
		return model.User{}, err
	}

	repoUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repoModel.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, model.ErrUserNotFound
	}
	if err != nil {
		return model.User{}, err
	}

	return converter.ToModelUser(repoUser)
}
