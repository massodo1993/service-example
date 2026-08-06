package v1

import (
	"github.com/massodo1993/service-example/order/internal/service"
)

type api struct {
	OrderService service.OrderService
}

func NewApi(orderService service.OrderService) *api {
	return &api{
		OrderService: orderService,
	}
}
