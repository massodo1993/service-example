package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/model"
	repoModel "github.com/massodo1993/service-example/order/internal/repository/model"
)

// ToRepoOrder конвертирует доменную модель заказа в модель репозитория.
func ToRepoOrder(order model.Order) repoModel.Order {
	partsUUIDs := make([]string, 0, len(order.PartsUUIDs))
	for _, partUUID := range order.PartsUUIDs {
		partsUUIDs = append(partsUUIDs, partUUID.String())
	}

	var transactionUUID *string
	if order.TransactionUUID != nil {
		txUUID := order.TransactionUUID.String()
		transactionUUID = &txUUID
	}

	var paymentMethod *repoModel.PaymentMethod
	if order.PaymentMethod != nil {
		pm := repoModel.PaymentMethod(*order.PaymentMethod)
		paymentMethod = &pm
	}

	return repoModel.Order{
		OrderUUID:       order.OrderUUID.String(),
		UserUUID:        order.UserUUID.String(),
		PartsUUIDs:      partsUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          repoModel.Status(order.Status),
	}
}

// ToModelOrder конвертирует модель репозитория заказа в доменную модель.
func ToModelOrder(order repoModel.Order) (model.Order, error) {
	orderUUID, err := uuid.Parse(order.OrderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("некорректный order_uuid %q: %w", order.OrderUUID, err)
	}

	userUUID, err := uuid.Parse(order.UserUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("некорректный user_uuid %q: %w", order.UserUUID, err)
	}

	partsUUIDs := make([]uuid.UUID, 0, len(order.PartsUUIDs))
	for _, partUUID := range order.PartsUUIDs {
		parsed, err := uuid.Parse(partUUID)
		if err != nil {
			return model.Order{}, fmt.Errorf("некорректный part_uuid %q: %w", partUUID, err)
		}

		partsUUIDs = append(partsUUIDs, parsed)
	}

	var transactionUUID *uuid.UUID
	if order.TransactionUUID != nil {
		txUUID, err := uuid.Parse(*order.TransactionUUID)
		if err != nil {
			return model.Order{}, fmt.Errorf("некорректный transaction_uuid %q: %w", *order.TransactionUUID, err)
		}

		transactionUUID = &txUUID
	}

	var paymentMethod *model.PaymentMethod
	if order.PaymentMethod != nil {
		pm := model.PaymentMethod(*order.PaymentMethod)
		paymentMethod = &pm
	}

	return model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartsUUIDs:      partsUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.Status(order.Status),
	}, nil
}
