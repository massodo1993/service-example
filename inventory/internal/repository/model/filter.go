package model

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []CategoryType
	ManufacturerCountries []string
	Tags                  []string
}
