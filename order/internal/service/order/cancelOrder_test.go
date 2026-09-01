package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	orderUUID := uuid.New()

	s.orderRepository.On("CancelOrder", s.ctx, orderUUID.String()).Return(nil)

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrderRepoError() {
	var (
		orderUUID = uuid.New()
		repoErr   = gofakeit.Error()
	)

	s.orderRepository.On("CancelOrder", s.ctx, orderUUID.String()).Return(repoErr)

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Require().ErrorIs(err, repoErr)
}
