package payment

import (
	def "github.com/massodo1993/service-example/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
