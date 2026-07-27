package model

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []CategoryType
	ManufacturerCountries []string
	Tags                  []string
}

func (f PartsFilter) IsEmpty() bool {
	return len(f.UUIDs) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Tags) == 0
}
