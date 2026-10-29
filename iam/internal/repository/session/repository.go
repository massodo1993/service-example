package session

import (
	"github.com/massodo1993/service-example/platform/pkg/cache"

	def "github.com/massodo1993/service-example/iam/internal/repository"
)

const (
	sessionKeyPrefix      = "iam:session:"
	userSessionsKeyPrefix = "iam:user_sessions:"
)

var _ def.SessionRepository = (*repository)(nil)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}
