package v1

import (
	"context"
	"errors"

	"github.com/massodo1993/service-example/order/internal/model"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.OrderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.CancelOrderNotFound{Code: "NOT_FOUND", Message: err.Error()}, nil
		} else if errors.Is(err, model.ErrOrderCancelConflict) {
			return &orderV1.CancelOrderConflict{Code: "CONFLICT", Message: err.Error()}, nil
		}

		return nil, err
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
