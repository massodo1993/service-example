package cmd

import (
	"sync"

	"github.com/google/uuid"
)

type OrderStatus int

const (
	UNKNOWN OrderStatus = iota
	CARD
	SBP
	CREDIT_CARD
	INVESTOR_MONEY
)

func (s OrderStatus) String() string {
	switch s {
	case UNKNOWN:
		return "UNKNOWN"
	case CARD:
		return "CARD"
	case SBP:
		return "SBP"
	case CREDIT_CARD:
		return "CREDIT_CARD"
	case INVESTOR_MONEY:
		return "INVESTOR_MONEY"
	default:
		return "UNKNOWN"
	}
}

type Order struct {
	orderUUID       uuid.UUID
	userUUID        uuid.UUID
	partsUUIDs      []uuid.UUID
	totalPrice      float64
	transactionUUID *uuid.UUID
	paymentMethod   *string
	status          OrderStatus
}

func NewOrder() *Order {
	return &Order{
		orderUUID: uuid.New(),
	}
}

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*Order),
	}
}

func (s *OrderStorage) GetOrder(uuid uuid.UUID) *Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, has := s.orders[uuid]
	if !has {
		return nil
	}

	return order
}

func (s *OrderStorage) SetOrder(order *Order) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.orders[order.orderUUID] = order
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

//func (h *OrderHandler) GetOrderByUUID(_ context.Context, params )

func main() {

}
