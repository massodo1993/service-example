package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/massodo1993/service-example/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		userUUID  = uuid.New()
		partUUID1 = uuid.New()
		partUUID2 = uuid.New()
		partsUUID = []uuid.UUID{partUUID1, partUUID2}

		parts = []model.Part{
			{PartUUID: partUUID1, Price: 100.5},
			{PartUUID: partUUID2, Price: 200},
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, partsUUID).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.AnythingOfType("model.Order")).Return(nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partsUUID)
	s.Require().NoError(err)
	s.Require().Equal(userUUID, order.UserUUID)
	s.Require().Equal(300.5, order.TotalPrice)
	s.Require().Equal(model.STATUS_PENDING_PAYMENT, order.Status)
}

func (s *ServiceSuite) TestCreateOrderPartNotFound() {
	var (
		userUUID  = uuid.New()
		partUUID  = uuid.New()
		partsUUID = []uuid.UUID{partUUID}
	)

	s.inventoryClient.On("ListParts", s.ctx, partsUUID).Return([]model.Part{}, nil)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partsUUID)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestCreateOrderInventoryError() {
	var (
		userUUID  = uuid.New()
		partsUUID = []uuid.UUID{uuid.New()}
		clientErr = gofakeit.Error()
	)

	s.inventoryClient.On("ListParts", s.ctx, partsUUID).Return(nil, clientErr)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partsUUID)
	s.Require().ErrorIs(err, clientErr)
	s.Require().Empty(order)
}

func (s *ServiceSuite) TestCreateOrderRepoError() {
	var (
		userUUID  = uuid.New()
		partUUID  = uuid.New()
		partsUUID = []uuid.UUID{partUUID}
		repoErr   = gofakeit.Error()

		parts = []model.Part{{PartUUID: partUUID, Price: 10}}
	)

	s.inventoryClient.On("ListParts", s.ctx, partsUUID).Return(parts, nil)
	s.orderRepository.On("CreateOrder", s.ctx, mock.AnythingOfType("model.Order")).Return(repoErr)

	order, err := s.service.CreateOrder(s.ctx, userUUID, partsUUID)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(order)
}
