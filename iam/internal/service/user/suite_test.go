package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/massodo1993/service-example/iam/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	userRepository *repoMocks.UserRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.userRepository = repoMocks.NewUserRepository(s.T())

	s.service = NewService(s.userRepository)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
