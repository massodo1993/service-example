package model

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
	OrderUUID       string         `db:"order_uuid"`
	UserUUID        string         `db:"user_uuid"`
	PartsUUIDs      []string       `db:"parts_uuids"`
	TotalPrice      float64        `db:"total_price"`
	TransactionUUID *string        `db:"transaction_uuid"`
	PaymentMethod   *PaymentMethod `db:"payment_method"`
	Status          Status         `db:"status"`
}
