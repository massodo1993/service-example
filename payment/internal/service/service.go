package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payment model.Payment) (uuid.UUID, error)
}
