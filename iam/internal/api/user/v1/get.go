package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/massodo1993/service-example/iam/internal/converter"
	userv1 "github.com/massodo1993/service-example/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "некорректный user_uuid")
	}

	user, err := a.userService.GetUser(ctx, userUUID)
	if err != nil {
		return nil, converter.ToGRPCError(err)
	}

	return &userv1.GetUserResponse{User: converter.ToProtoUser(user)}, nil
}
