package model

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType int

const (
	CATEGORY_TYPE_UNKNOW CategoryType = iota
	CATEGORY_TYPE_ENGINE
	CATEGORY_TYPE_FUEL
	CATEGORY_TYPE_PORTHOLE
	CATEGORY_TYPE_WING
)

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

type Value struct {
	String_value *string
	Int64_value  *int64
	Double_value *float64
	Bool_value   *bool
}

type Part struct {
	Uuid          uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	CategoryType  CategoryType
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
}
