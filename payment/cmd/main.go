package main

import "github.com/google/uuid"

type Transaction struct {
	order_uuid     uuid.UUID
	user_uuid      uuid.UUID
	payment_method PaymentMethod
}

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
	case PM_CARD:
		return "Банковская карта"
	case PM_SBP:
		return "SBP"
	case PM_CREDIT_CARD:
		return "Система быстрых платежей"
	case PM_INVESTOR_MONEY:
		return "Кредитная карта"
	default:
		return "Неизвестный способ"
	}
}

func main() {

}
