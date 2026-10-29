package v1

import (
	"github.com/massodo1993/service-example/iam/internal/service"
	userv1 "github.com/massodo1993/service-example/shared/pkg/proto/user/v1"
)

type api struct {
	userv1.UnimplementedUserServiceServer

	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
