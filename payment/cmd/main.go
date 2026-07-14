package main

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	paymentv1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

type payemntService struct {
	paymentv1.UnimplementedPayemntServiceServer

	mu       sync.Mutex
	payments map[string]*paymentv1.PaymentMethod
}

func (ps *payemntService) PayOrder(_ context.Context, request *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	uuid := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", uuid)
	return &paymentv1.PayOrderResponse{TransactionUuid: uuid.String()}, nil
}

func main() {

}
