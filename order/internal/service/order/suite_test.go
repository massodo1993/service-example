package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/massodo1993/service-example/order/internal/client/mocks"
	repoMocks "github.com/massodo1993/service-example/order/internal/repository/mocks"
	serviceMocks "github.com/massodo1993/service-example/order/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *repoMocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient
	orderProducer   *serviceMocks.OrderProducerService

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	s.orderProducer = serviceMocks.NewOrderProducerService(s.T())

	s.service = NewService(s.orderRepository, s.inventoryClient, s.paymentClient, s.orderProducer)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
