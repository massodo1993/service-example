package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/massodo1993/service-example/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	partRepository *repoMocks.PartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.partRepository = repoMocks.NewPartRepository(s.T())
	s.service = NewService(s.partRepository)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
