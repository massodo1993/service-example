package v1

import (
	"context"

	"github.com/massodo1993/service-example/iam/internal/converter"
	authv1 "github.com/massodo1993/service-example/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	sessionUUID, err := a.authService.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, converter.ToGRPCError(err)
	}

	return &authv1.LoginResponse{SessionUuid: sessionUUID.String()}, nil
}
