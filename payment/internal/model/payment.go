package model

import "github.com/google/uuid"

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

type Payment struct {
	OrderUUID     uuid.UUID
	UserUUID      uuid.UUID
	PaymentMethod PaymentMethod
}
