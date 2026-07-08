package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
)

type PaymentMethod int

const (
	PM_UNKNOWN PaymentMethod = iota
	PM_CARD
	PM_SBP
	PM_CREDIT_CARD
	PM_INVESTOR_MONEY
)

func (pm PaymentMethod) String() string {
	switch pm {
	case PM_UNKNOWN:
		return "UNKNOWN"
	case PM_CARD:
		return "CARD"
	case PM_SBP:
		return "SBP"
	case PM_CREDIT_CARD:
		return "CREDIT_CARD"
	case PM_INVESTOR_MONEY:
		return "INVESTOR_MONEY"
	default:
		return "UNKNOWN"
	}
}

func (pm *PaymentMethod) fromString(value string) error {
	switch value {
	case "UNKNOWN":
		*pm = PM_UNKNOWN
	case "CARD":
		*pm = PM_CARD
	case "SBP":
		*pm = PM_SBP
	case "CREDIT_CARD":
		*pm = PM_CREDIT_CARD
	case "INVESTOR_MONEY":
		*pm = PM_INVESTOR_MONEY
	default:
		return fmt.Errorf("error case")

	}
	return nil
}

type Status int

const (
	S_UNKNOWN Status = iota
	S_PENDING_PAYMENT
	S_PAID
	S_CANCELLED
)

func (s Status) String() string {
	switch s {
	case S_UNKNOWN:
		return "UNKNOWN"
	case S_PENDING_PAYMENT:
		return "PENDING_PAYMENT"
	case S_PAID:
		return "PAID"
	case S_CANCELLED:
		return "CANCELLED"
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
	paymentMethod   *PaymentMethod
	status          Status
}

func NewOrder(id uuid.UUID, parts []uuid.UUID, status Status) *Order {
	return &Order{
		orderUUID:  uuid.New(),
		userUUID:   id,
		partsUUIDs: parts,
		status:     status,
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

func (h *OrderHandler) GetOrderByUuid(_ context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.GetOrderByUuidNotFound{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("Order with %s not found", params.OrderUUID),
		}, nil
	}

	return &orderV1.OrderDto{
		OrderUUID:  order.orderUUID,
		UserUUID:   order.userUUID,
		PartUuids:  order.partsUUIDs,
		TotalPrice: order.totalPrice,
		Status:     orderV1.OrderStatus(order.status.String()),
	}, nil
}

func (h *OrderHandler) CreateOrder(_ context.Context, request *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	//TO-DO проверка деталей и посчитать сумму
	order := NewOrder(request.GetUserUUID(), request.GetPartUuids(), S_PENDING_PAYMENT)
	h.storage.SetOrder(order)
	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.orderUUID,
		TotalPrice: order.totalPrice,
	}, nil
}

func (h *OrderHandler) PayOrder(_ context.Context, request *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.PayOrderNotFound{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("Order with %s not found", params.OrderUUID),
		}, nil
	}
	//TO-DO оплата
	order.paymentMethod.fromString(string(request.PaymentMethod))
	order.status = S_PAID
	return &orderV1.PayOrderResponse{
		TransactionUUID: uuid.New(), //TO-do поменяй на нормальный uuid
	}, nil
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.CancelOrderNotFound{
			Code:    "NOT_FOUND",
			Message: fmt.Sprintf("Order with %s not found", params.OrderUUID),
		}, nil
	}

	if order.status == S_PAID {
		return &orderV1.CancelOrderConflict{
			Code:    "CONFLICT",
			Message: "Payed order can't be canceled",
		}, nil
	}

	if order.status == S_PENDING_PAYMENT {
		return &orderV1.CancelOrderNoContent{}, nil
	}

	return &orderV1.CancelOrderInternalServerError{}, nil
}

func main() {
	storage := NewOrderStorage()
	orderHandler := NewOrderHandler(storage)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("OpenAPI server creation error: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	router.Mount("/", orderServer)
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", "8001"),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Server start")
		err = server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Server stop")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	log.Printf("Server stop")
}
