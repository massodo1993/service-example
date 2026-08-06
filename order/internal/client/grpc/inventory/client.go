package inventory

import (
	def "github.com/massodo1993/service-example/order/internal/client"
	inventoryV1 "github.com/massodo1993/service-example/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient inventoryV1.InventoryServiceClient
}

func NewClient(generatedClient inventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
