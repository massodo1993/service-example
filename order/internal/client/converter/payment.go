package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
	paymentV1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

func ToProtoPaymentMethod(paymentMethod model.PaymentMethod) paymentV1.PaymentMethod {
	switch paymentMethod {
	case model.PM_CARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PM_SBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PM_CREDIT_CARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PM_INVESTOR_MONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func ToModelTransactionUUID(response *paymentV1.PayOrderResponse) (uuid.UUID, error) {
	transactionUUID, err := uuid.Parse(response.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("некорректный transaction_uuid %q: %w", response.GetTransactionUuid(), err)
	}

	return transactionUUID, nil
}
