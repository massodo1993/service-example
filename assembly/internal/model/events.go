package model

// OrderPaidEvent — входящее событие "заказ оплачен".
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

// ShipAssembledEvent — исходящее событие "корабль собран".
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
