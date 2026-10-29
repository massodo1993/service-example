package user

import (
	"github.com/massodo1993/service-example/iam/internal/repository"
	def "github.com/massodo1993/service-example/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
