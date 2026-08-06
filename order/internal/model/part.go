package model

import (
	"time"

	"github.com/google/uuid"
)

type Category int

const (
	CATEGORY_UNKNOWN Category = iota
	CATEGORY_ENGINE
	CATEGORY_FUEL
	CATEGORY_PORTHOLE
	CATEGORY_WING
)

func (c Category) String() string {
	switch c {
	case CATEGORY_UNKNOWN:
		return "UNKNOWN"
	case CATEGORY_ENGINE:
		return "ENGINE"
	case CATEGORY_FUEL:
		return "FUEL"
	case CATEGORY_PORTHOLE:
		return "PORTHOLE"
	case CATEGORY_WING:
		return "WING"
	default:
		return "UNKNOWN"
	}
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// Value — доменное представление oneof-поля metadata из inventory.
// Заполнено ровно одно поле, остальные nil.
type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Part struct {
	PartUUID      uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
