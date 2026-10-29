package auth

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *ServiceSuite) hashedUser(login, password string) model.User {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	s.Require().NoError(err)

	return model.User{UUID: uuid.New(), Info: model.UserInfo{Login: login}, PasswordHash: string(hash)}
}

func (s *ServiceSuite) TestLoginSuccess() {
	var (
		login    = gofakeit.Username()
		password = gofakeit.Password(true, true, true, true, false, 12)
		user     = s.hashedUser(login, password)
	)

	s.userRepository.On("GetByLogin", s.ctx, login).Return(user, nil)
	s.sessionRepository.On("Create", s.ctx, mock.AnythingOfType("uuid.UUID"), user.UUID, sessionTTL).Return(nil)

	got, err := s.service.Login(s.ctx, login, password)
	s.Require().NoError(err)
	s.Require().NotEqual(uuid.Nil, got)
}

func (s *ServiceSuite) TestLoginWrongPassword() {
	login := gofakeit.Username()
	user := s.hashedUser(login, "correct-password")

	s.userRepository.On("GetByLogin", s.ctx, login).Return(user, nil)

	got, err := s.service.Login(s.ctx, login, "wrong-password")
	s.Require().ErrorIs(err, model.ErrInvalidCredentials)
	s.Require().Equal(uuid.Nil, got)
}

func (s *ServiceSuite) TestLoginUserNotFound() {
	login := gofakeit.Username()

	s.userRepository.On("GetByLogin", s.ctx, login).Return(model.User{}, model.ErrUserNotFound)

	got, err := s.service.Login(s.ctx, login, "password")
	s.Require().ErrorIs(err, model.ErrInvalidCredentials)
	s.Require().Equal(uuid.Nil, got)
}

func (s *ServiceSuite) TestLoginSessionRepoError() {
	var (
		login    = gofakeit.Username()
		password = "password"
		user     = s.hashedUser(login, password)
		repoErr  = gofakeit.Error()
	)

	s.userRepository.On("GetByLogin", s.ctx, login).Return(user, nil)
	s.sessionRepository.On("Create", s.ctx, mock.AnythingOfType("uuid.UUID"), user.UUID, sessionTTL).Return(repoErr)

	got, err := s.service.Login(s.ctx, login, password)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Equal(uuid.Nil, got)
}
