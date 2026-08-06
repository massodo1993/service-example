package converter

import (
	"github.com/massodo1993/service-example/order/internal/model"
	orderV1 "github.com/massodo1993/service-example/shared/pkg/openapi/order/v1"
)

// ToOrderDto конвертирует доменную модель заказа в модель API.
func ToOrderDto(order model.Order) orderV1.OrderDto {
	dto := orderV1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartsUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     ToOrderStatus(order.Status),
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptUUID(*order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		dto.PaymentMethod = orderV1.NewOptPaymentMethod(ToPaymentMethodDto(*order.PaymentMethod))
	}

	return dto
}

// ToOrderStatus конвертирует доменный статус заказа в статус API.
func ToOrderStatus(status model.Status) orderV1.OrderStatus {
	switch status {
	case model.STATUS_PENDING_PAYMENT:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.STATUS_PAID:
		return orderV1.OrderStatusPAID
	case model.STATUS_CANCELLED:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

// ToPaymentMethodDto конвертирует доменный способ оплаты в способ оплаты API.
func ToPaymentMethodDto(paymentMethod model.PaymentMethod) orderV1.PaymentMethod {
	switch paymentMethod {
	case model.PM_CARD:
		return orderV1.PaymentMethodCARD
	case model.PM_SBP:
		return orderV1.PaymentMethodSBP
	case model.PM_CREDIT_CARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.PM_INVESTOR_MONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

// ToModelPaymentMethod конвертирует способ оплаты из запроса API в доменный.
func ToModelPaymentMethod(paymentMethod orderV1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case orderV1.PaymentMethodCARD:
		return model.PM_CARD
	case orderV1.PaymentMethodSBP:
		return model.PM_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PM_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PM_INVESTOR_MONEY
	default:
		return model.PM_UNKNOWN
	}
}
