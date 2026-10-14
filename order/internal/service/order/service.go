package order

import (
	"github.com/massodo1993/service-example/order/internal/client"
	"github.com/massodo1993/service-example/order/internal/repository"
	def "github.com/massodo1993/service-example/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
	orderProducer   def.OrderProducerService
}

func NewService(
	orderRepository repository.OrderRepository,
	initInventoryClient client.InventoryClient,
	initPaymentClient client.PaymentClient,
	orderProducer def.OrderProducerService,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: initInventoryClient,
		paymentClient:   initPaymentClient,
		orderProducer:   orderProducer,
	}
}
