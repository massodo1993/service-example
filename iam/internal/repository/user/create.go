package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/massodo1993/service-example/iam/internal/model"
	"github.com/massodo1993/service-example/iam/internal/repository/converter"
	repoModel "github.com/massodo1993/service-example/iam/internal/repository/model"
)

// pgUniqueViolation — код ошибки PostgreSQL при нарушении уникального ограничения.
const pgUniqueViolation = "23505"

func (r *repository) Create(ctx context.Context, info model.UserInfo, passwordHash string) (model.User, error) {
	methods, err := converter.ToRepoNotificationMethods(info.NotificationMethods)
	if err != nil {
		return model.User{}, err
	}

	rows, err := r.connect.Query(ctx,
		"INSERT INTO users (login, email, password_hash, notification_methods) VALUES ($1, $2, $3, $4) RETURNING uuid, login, email, password_hash, notification_methods, created_at, updated_at",
		info.Login, info.Email, passwordHash, methods,
	)
	if err != nil {
		return model.User{}, err
	}

	repoUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repoModel.User])
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return model.User{}, model.ErrUserAlreadyExists
		}

		return model.User{}, err
	}

	return converter.ToModelUser(repoUser)
}
