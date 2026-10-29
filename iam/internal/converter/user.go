package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/massodo1993/service-example/iam/internal/model"
	commonv1 "github.com/massodo1993/service-example/shared/pkg/proto/common/v1"
)

// ToProtoUser конвертирует доменную модель пользователя в proto (без хэша пароля).
func ToProtoUser(user model.User) *commonv1.User {
	methods := make([]*commonv1.NotificationMethod, 0, len(user.Info.NotificationMethods))
	for _, method := range user.Info.NotificationMethods {
		methods = append(methods, &commonv1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return &commonv1.User{
		Uuid: user.UUID.String(),
		Info: &commonv1.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: methods,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToModelUserInfo конвертирует proto-информацию о пользователе в доменную модель.
func ToModelUserInfo(info *commonv1.UserInfo) model.UserInfo {
	methods := make([]model.NotificationMethod, 0, len(info.GetNotificationMethods()))
	for _, method := range info.GetNotificationMethods() {
		methods = append(methods, model.NotificationMethod{
			ProviderName: method.GetProviderName(),
			Target:       method.GetTarget(),
		})
	}

	return model.UserInfo{
		Login:               info.GetLogin(),
		Email:               info.GetEmail(),
		NotificationMethods: methods,
	}
}
