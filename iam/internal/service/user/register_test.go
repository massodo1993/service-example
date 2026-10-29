package user

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/massodo1993/service-example/iam/internal/model"
)

func (s *ServiceSuite) TestRegisterSuccess() {
	var (
		userUUID = uuid.New()
		password = gofakeit.Password(true, true, true, true, false, 12)
		info     = model.UserInfo{Login: gofakeit.Username(), Email: gofakeit.Email()}
	)

	s.userRepository.On("Create", s.ctx, info, mock.MatchedBy(func(hash string) bool {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
	})).Return(model.User{UUID: userUUID, Info: info}, nil)

	got, err := s.service.Register(s.ctx, model.UserRegistrationInfo{Info: info, Password: password})
	s.Require().NoError(err)
	s.Require().Equal(userUUID, got)
}

func (s *ServiceSuite) TestRegisterInvalidArgument() {
	got, err := s.service.Register(s.ctx, model.UserRegistrationInfo{
		Info:     model.UserInfo{Login: gofakeit.Username()},
		Password: gofakeit.Password(true, true, true, true, false, 12),
	})
	s.Require().ErrorIs(err, model.ErrInvalidArgument)
	s.Require().Equal(uuid.Nil, got)
}

func (s *ServiceSuite) TestRegisterRepoError() {
	var (
		repoErr = model.ErrUserAlreadyExists
		info    = model.UserInfo{Login: gofakeit.Username(), Email: gofakeit.Email()}
	)

	s.userRepository.On("Create", s.ctx, info, mock.Anything).Return(model.User{}, repoErr)

	got, err := s.service.Register(s.ctx, model.UserRegistrationInfo{Info: info, Password: "password"})
	s.Require().ErrorIs(err, repoErr)
	s.Require().Equal(uuid.Nil, got)
}
