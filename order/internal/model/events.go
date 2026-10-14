package model

import "github.com/google/uuid"

// OrderPaidEvent — событие "заказ оплачен", публикуется в Kafka.
type OrderPaidEvent struct {
	EventUUID       uuid.UUID
	OrderUUID       uuid.UUID
	UserUUID        uuid.UUID
	PaymentMethod   string
	TransactionUUID uuid.UUID
}

// ShipAssembledEvent — событие "корабль собран", приходит из Kafka.
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
