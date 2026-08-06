package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/massodo1993/service-example/order/internal/model"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order, err := a.OrderService.CreateOrder(ctx, req.GetUserUUID(), req.GetPartUuids())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &orderV1.CreateOrderBadRequest{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			}, nil
		}

		return nil, fmt.Errorf("создание заказа: %w", err)
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
