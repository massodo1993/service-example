package user

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *ServiceSuite) TestGetUserSuccess() {
	var (
		userUUID = uuid.New()
		user     = model.User{UUID: userUUID, Info: model.UserInfo{Login: gofakeit.Username()}}
	)

	s.userRepository.On("GetByUUID", s.ctx, userUUID).Return(user, nil)

	got, err := s.service.GetUser(s.ctx, userUUID)
	s.Require().NoError(err)
	s.Require().Equal(user, got)
}

func (s *ServiceSuite) TestGetUserNotFound() {
	userUUID := uuid.New()

	s.userRepository.On("GetByUUID", s.ctx, userUUID).Return(model.User{}, model.ErrUserNotFound)

	got, err := s.service.GetUser(s.ctx, userUUID)
	s.Require().ErrorIs(err, model.ErrUserNotFound)
	s.Require().Equal(model.User{}, got)
}
