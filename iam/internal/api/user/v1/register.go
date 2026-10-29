package v1

import (
	"context"

	"github.com/massodo1993/service-example/iam/internal/converter"
	"github.com/massodo1993/service-example/iam/internal/model"
	userv1 "github.com/massodo1993/service-example/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	userUUID, err := a.userService.Register(ctx, model.UserRegistrationInfo{
		Info:     converter.ToModelUserInfo(req.GetInfo().GetInfo()),
		Password: req.GetInfo().GetPassword(),
	})
	if err != nil {
		return nil, converter.ToGRPCError(err)
	}

	return &userv1.RegisterResponse{UserUuid: userUUID.String()}, nil
}
