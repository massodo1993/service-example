package model

import "github.com/google/uuid"

type Status int

const (
	STATUS_UNKNOWN Status = iota
	STATUS_PENDING_PAYMENT
	STATUS_PAID
	STATUS_CANCELLED
)

func (s Status) String() string {
	switch s {
	case STATUS_UNKNOWN:
		return "UNKNOWN"
	case STATUS_PENDING_PAYMENT:
		return "PENDING_PAYMENT"
	case STATUS_PAID:
		return "PAID"
	case STATUS_CANCELLED:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
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

type Order struct {
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PartsUUIDs      []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          Status
}

func NewOrder(id uuid.UUID, parts []uuid.UUID, status Status, price float64) Order {
	return Order{
		OrderUUID:  uuid.New(),
		UserUUID:   id,
		PartsUUIDs: parts,
		Status:     status,
		TotalPrice: price,
	}
}
