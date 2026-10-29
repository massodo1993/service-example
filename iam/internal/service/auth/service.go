package auth

import (
	"time"

	"github.com/massodo1993/service-example/iam/internal/repository"
	def "github.com/massodo1993/service-example/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	sessionTTL        time.Duration
}

func NewService(
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
	sessionTTL time.Duration,
) *service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		sessionTTL:        sessionTTL,
	}
}
