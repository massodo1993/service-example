package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/payment/internal/model"
	paymentv1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

// ToDomainPayment конвертирует gRPC-запрос в доменную команду на оплату.
func ToDomainPayment(request *paymentv1.PayOrderRequest) (model.Payment, error) {
	orderUUID, err := uuid.Parse(request.GetOrderUuid())
	if err != nil {
		return model.Payment{}, fmt.Errorf("некорректный order_uuid %q: %w", request.GetOrderUuid(), err)
	}

	userUUID, err := uuid.Parse(request.GetUserUuid())
	if err != nil {
		return model.Payment{}, fmt.Errorf("некорректный user_uuid %q: %w", request.GetUserUuid(), err)
	}

	return model.Payment{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PaymentMethod: ToDomainPaymentMethod(request.GetPaymentMethod()),
	}, nil
}

// ToDomainPaymentMethod конвертирует способ оплаты из proto в доменный.
func ToDomainPaymentMethod(paymentMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PM_CARD
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PM_SBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PM_CREDIT_CARD
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PM_INVESTOR_MONEY
	default:
		return model.PM_UNKNOWN
	}
}
