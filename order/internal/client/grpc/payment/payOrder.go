package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/order/internal/client/converter"
	"github.com/massodo1993/service-example/order/internal/model"
	paymentV1 "github.com/massodo1993/service-example/shared/pkg/proto/payment/v1"
)

// PayOrder проводит оплату заказа через payment и возвращает UUID транзакции.
func (c *client) PayOrder(ctx context.Context, orderUuid, userUuid uuid.UUID, pm model.PaymentMethod) (uuid.UUID, error) {
	response, err := c.generatedClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUuid.String(),
		UserUuid:      userUuid.String(),
		PaymentMethod: converter.ToProtoPaymentMethod(pm),
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплата заказа в payment: %w", err)
	}

	transactionUUID, err := converter.ToModelTransactionUUID(response)
	if err != nil {
		return uuid.Nil, fmt.Errorf("конвертация ответа payment: %w", err)
	}

	return transactionUUID, nil
}
