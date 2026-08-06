package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/massodo1993/service-example/payment/internal/converter"
	"github.com/massodo1993/service-example/payment/internal/model"
	paymentv1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, request *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	payment, err := converter.ToDomainPayment(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transactionUUID, err := a.paymentService.PayOrder(ctx, payment)
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
