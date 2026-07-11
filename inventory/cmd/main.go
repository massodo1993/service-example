package main

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Part struct {
	uuid           uuid.UUID
	name           string
	description    string
	price          float32
	stock_quantity int64
	category       Category
	dimensions     Dimensions
	manufacturer   Manufacturer
	tags           []string
	metadata       map[string]Value
	created_at     time.Time
	updated_at     time.Time
}

type Category int

const (
	C_UNKNOWN Category = iota
	C_ENGINE
	C_FUEL
	C_PORTHOLE
	C_WING
)

func (c Category) String() string {
	switch c {
	case 1:
		return "UNKNOWN"
	case 2:
		return "ENGINE"
	case 3:
		return "FUEL"
	case 4:
		return "PORTHOLE"
	case 5:
		return "WING"
	default:
		return "UNKNOWN"
	}
}

type Dimensions struct {
	length float64
	width  float64
	height float64
	weight float64
}

type Manufacturer struct {
	name    string
	country string
	website string
}

type Value struct {
	string_value string
	int64_value  int64
	double_value float64
	bool_value   bool
}

func NewPart() Part {
	return Part{
		uuid: uuid.New(),
	}
}

type PartStorage struct {
	mu   sync.RWMutex
	part map[int]*Part
}

func NewPartStorage() *PartStorage {
	return &PartStorage{
		part: make(map[int]*Part),
	}
}

func main() {

}
