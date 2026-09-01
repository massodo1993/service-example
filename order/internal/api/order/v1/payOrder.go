package v1

import (
	"context"
	"errors"

	"github.com/massodo1993/service-example/order/internal/converter"
	"github.com/massodo1993/service-example/order/internal/model"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, request *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	paymentMethod := converter.ToModelPaymentMethod(request.GetPaymentMethod())
	transactionUuid, err := a.OrderService.PayOrder(ctx, params.OrderUUID, paymentMethod)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.PayOrderNotFound{Code: "NOT_FOUND", Message: err.Error()}, nil
		} else if errors.Is(err, model.ErrOrderPayConflict) {
			return &orderV1.PayOrderBadGateway{Code: "BAD_GETWAY", Message: err.Error()}, nil
		}

		return nil, err
	}

	return &orderV1.PayOrderResponse{TransactionUUID: transactionUuid}, nil
}
