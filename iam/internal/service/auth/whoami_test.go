package auth

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *ServiceSuite) TestWhoamiSuccess() {
	var (
		sessionUUID = uuid.New()
		userUUID    = uuid.New()
		user        = model.User{UUID: userUUID, Info: model.UserInfo{Login: gofakeit.Username()}}
	)

	s.sessionRepository.On("Get", s.ctx, sessionUUID).Return(userUUID, nil)
	s.userRepository.On("GetByUUID", s.ctx, userUUID).Return(user, nil)

	got, err := s.service.Whoami(s.ctx, sessionUUID)
	s.Require().NoError(err)
	s.Require().Equal(user, got)
}

func (s *ServiceSuite) TestWhoamiSessionNotFound() {
	sessionUUID := uuid.New()

	s.sessionRepository.On("Get", s.ctx, sessionUUID).Return(uuid.Nil, model.ErrSessionNotFound)

	got, err := s.service.Whoami(s.ctx, sessionUUID)
	s.Require().ErrorIs(err, model.ErrSessionNotFound)
	s.Require().Equal(model.User{}, got)
}

func (s *ServiceSuite) TestWhoamiUserRepoError() {
	var (
		sessionUUID = uuid.New()
		userUUID    = uuid.New()
		repoErr     = gofakeit.Error()
	)

	s.sessionRepository.On("Get", s.ctx, sessionUUID).Return(userUUID, nil)
	s.userRepository.On("GetByUUID", s.ctx, userUUID).Return(model.User{}, repoErr)

	got, err := s.service.Whoami(s.ctx, sessionUUID)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Equal(model.User{}, got)
}
