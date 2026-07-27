package part

import (
	"github.com/massodo1993/service-example/inventory/internal/repository"
	def "github.com/massodo1993/service-example/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	PartRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *service {
	return &service{
		PartRepository: partRepository,
	}
}
