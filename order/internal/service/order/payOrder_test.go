package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		transactionUUID = uuid.New()
		paymentMethod   = model.PM_CARD

		order = model.Order{OrderUUID: orderUUID, UserUUID: userUUID, Status: model.STATUS_PENDING_PAYMENT}
	)

	s.orderRepository.On("GetOrderByUuid", s.ctx, orderUUID.String()).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil)
	s.orderRepository.On("PayOrder", s.ctx, orderUUID.String(), transactionUUID.String(), paymentMethod).Return(nil)

	got, err := s.service.PayOrder(s.ctx, orderUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, got)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	var (
		orderUUID = uuid.New()
		repoErr   = model.ErrOrderNotFound
	)

	s.orderRepository.On("GetOrderByUuid", s.ctx, orderUUID.String()).Return(model.Order{}, repoErr)

	got, err := s.service.PayOrder(s.ctx, orderUUID, model.PM_CARD)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Equal(uuid.Nil, got)
}

func (s *ServiceSuite) TestPayOrderPaymentClientError() {
	var (
		orderUUID = uuid.New()
		userUUID  = uuid.New()
		clientErr = gofakeit.Error()

		order = model.Order{OrderUUID: orderUUID, UserUUID: userUUID, Status: model.STATUS_PENDING_PAYMENT}
	)

	s.orderRepository.On("GetOrderByUuid", s.ctx, orderUUID.String()).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, model.PM_CARD).Return(uuid.Nil, clientErr)

	got, err := s.service.PayOrder(s.ctx, orderUUID, model.PM_CARD)
	s.Require().ErrorIs(err, clientErr)
	s.Require().Equal(uuid.Nil, got)
}

func (s *ServiceSuite) TestPayOrderRepoError() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		transactionUUID = uuid.New()
		repoErr         = gofakeit.Error()

		order = model.Order{OrderUUID: orderUUID, UserUUID: userUUID, Status: model.STATUS_PENDING_PAYMENT}
	)

	s.orderRepository.On("GetOrderByUuid", s.ctx, orderUUID.String()).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, orderUUID, userUUID, model.PM_CARD).Return(transactionUUID, nil)
	s.orderRepository.On("PayOrder", s.ctx, orderUUID.String(), transactionUUID.String(), model.PM_CARD).Return(repoErr)

	got, err := s.service.PayOrder(s.ctx, orderUUID, model.PM_CARD)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Equal(uuid.Nil, got)
}
