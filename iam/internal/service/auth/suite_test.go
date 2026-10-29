package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/massodo1993/service-example/iam/internal/repository/mocks"
)

const sessionTTL = 24 * time.Hour

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	userRepository    *repoMocks.UserRepository
	sessionRepository *repoMocks.SessionRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.userRepository = repoMocks.NewUserRepository(s.T())
	s.sessionRepository = repoMocks.NewSessionRepository(s.T())

	s.service = NewService(s.userRepository, s.sessionRepository, sessionTTL)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
