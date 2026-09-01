package model

import "errors"

var (
	ErrOrderNotFound       = errors.New("order with this uuid doesn't exist")
	ErrOrderPayConflict    = errors.New("order is already paid or cancelled")
	ErrOrderCancelConflict = errors.New("order is already paid or cancelled")
	ErrPartNotFound        = errors.New("wrong parts uuid")
)
