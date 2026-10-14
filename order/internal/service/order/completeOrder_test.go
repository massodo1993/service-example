package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestCompleteOrderSuccess() {
	orderUUID := uuid.New()

	s.orderRepository.On("CompleteOrder", s.ctx, orderUUID.String()).Return(nil)

	err := s.service.CompleteOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCompleteOrderRepoError() {
	var (
		orderUUID = uuid.New()
		repoErr   = gofakeit.Error()
	)

	s.orderRepository.On("CompleteOrder", s.ctx, orderUUID.String()).Return(repoErr)

	err := s.service.CompleteOrder(s.ctx, orderUUID)
	s.Require().ErrorIs(err, repoErr)
}
