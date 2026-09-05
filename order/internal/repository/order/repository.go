package order

import (
	"github.com/jackc/pgx/v5/pgxpool"
	def "github.com/massodo1993/service-example/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	connect *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		connect: pool,
	}
}
