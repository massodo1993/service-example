package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/massodo1993/service-example/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	connect *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		connect: pool,
	}
}
