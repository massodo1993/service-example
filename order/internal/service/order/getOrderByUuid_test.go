package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderByUuidSuccess() {
	var (
		orderUUID = uuid.New()
		expected  = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   uuid.New(),
			TotalPrice: 42,
			Status:     model.STATUS_PENDING_PAYMENT,
		}
	)

	s.orderRepository.On("GetOrderByUuid", s.ctx, orderUUID.String()).Return(expected, nil)

	order, err := s.service.GetOrderByUuid(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(expected, order)
}

func (s *ServiceSuite) TestGetOrderByUuidRepoError() {
	var (
		orderUUID = uuid.New()
		repoErr   = gofakeit.Error()
	)

	s.orderRepository.On("GetOrderByUuid", s.ctx, orderUUID.String()).Return(model.Order{}, repoErr)

	order, err := s.service.GetOrderByUuid(s.ctx, orderUUID)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(order)
}
