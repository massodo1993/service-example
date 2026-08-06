package payment

import (
	def "github.com/massodo1993/service-example/order/internal/client"
	paymentV1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(generatedClient paymentV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
