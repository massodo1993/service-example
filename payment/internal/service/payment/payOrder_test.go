package payment

import (
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	payment := model.Payment{
		OrderUUID:     uuid.New(),
		UserUUID:      uuid.New(),
		PaymentMethod: model.PM_CARD,
	}

	transactionUUID, err := s.service.PayOrder(s.ctx, payment)
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, transactionUUID)
}

func (s *ServiceSuite) TestPayOrderUnknownMethod() {
	payment := model.Payment{
		OrderUUID:     uuid.New(),
		UserUUID:      uuid.New(),
		PaymentMethod: model.PM_UNKNOWN,
	}

	transactionUUID, err := s.service.PayOrder(s.ctx, payment)
	s.Require().ErrorIs(err, model.ErrInvalidPaymentMethod)
	s.Require().Equal(uuid.Nil, transactionUUID)
}
