package app

import (
	"context"

	payemntV1API "github.com/massodo1993/service-example/payment/internal/api/payment/v1"
	"github.com/massodo1993/service-example/payment/internal/service"
	paymentService "github.com/massodo1993/service-example/payment/internal/service/payment"
	paymentv1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentv1API   paymentv1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.paymentv1API == nil {
		d.paymentv1API = payemntV1API.NewAPI(d.PartService(ctx))
	}

	return d.paymentv1API
}

func (d *diContainer) PartService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
