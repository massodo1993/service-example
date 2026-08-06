package order

import (
	"sync"

	def "github.com/massodo1993/service-example/order/internal/repository"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]*repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]*repoModel.Order),
	}
}
