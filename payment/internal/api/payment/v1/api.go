package v1

import (
	"github.com/massodo1993/service-example/payment/internal/service"
	paymentv1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
