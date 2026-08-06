package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/payment/internal/model"
)

func (s *service) PayOrder(_ context.Context, payment model.Payment) (uuid.UUID, error) {
	if payment.PaymentMethod == model.PM_UNKNOWN {
		return uuid.Nil, model.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return transactionUUID, nil
}
