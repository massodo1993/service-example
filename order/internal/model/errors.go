package model

import "errors"

var ErrOrderNotFound = errors.New("order with this uuid doesn't exist")
var ErrOrderPayConflict = errors.New("order is already paid or cancelled")
var ErrOrderCancelConflict = errors.New("order is already paid or cancelled")
var ErrPartNotFound = errors.New("wrong parts uuid")
