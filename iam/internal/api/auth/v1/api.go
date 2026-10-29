package v1

import (
	"github.com/massodo1993/service-example/iam/internal/service"
	authv1 "github.com/massodo1993/service-example/shared/pkg/proto/auth/v1"
)

type api struct {
	authv1.UnimplementedAuthServiceServer

	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
