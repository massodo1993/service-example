package v1

import (
	"context"
	"errors"

	"github.com/massodo1993/service-example/order/internal/converter"
	"github.com/massodo1993/service-example/order/internal/model"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order, err := a.OrderService.GetOrderByUuid(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.GetOrderByUuidNotFound{Code: "NOT_FOUND", Message: err.Error()}, nil
		}

		return nil, err
	}

	dto := converter.ToOrderDto(order)

	return &dto, nil
}
