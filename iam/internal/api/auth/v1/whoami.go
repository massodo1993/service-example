package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/massodo1993/service-example/iam/internal/converter"
	authv1 "github.com/massodo1993/service-example/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	sessionUUID, err := uuid.Parse(req.GetSessionUuid())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "некорректный session_uuid")
	}

	user, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		return nil, converter.ToGRPCError(err)
	}

	return &authv1.WhoamiResponse{User: converter.ToProtoUser(user)}, nil
}
