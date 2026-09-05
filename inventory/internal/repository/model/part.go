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
	Uuid          uuid.UUID        `bson:"part_uuid"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description"`
	Price         float64          `bson:"price"`
	StockQuantity int              `bson:"stock_quantity"`
	CategoryType  CategoryType     `bson:"category_type"`
	Dimensions    *Dimensions      `bson:"dimensions"`
	Manufacturer  *Manufacturer    `bson:"manufacturer"`
	Tags          []string         `bson:"tags"`
	Metadata      map[string]Value `bson:"metadata"`
	CreatedAt     time.Time        `bson:"created_at"`
	UpdatedAt     *time.Time       `bson:"updated_at"`
	DeletedAt     *time.Time       `bson:"deleted_at"`
}
