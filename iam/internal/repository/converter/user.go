package converter

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/massodo1993/service-example/iam/internal/model"
	repoModel "github.com/massodo1993/service-example/iam/internal/repository/model"
)

// ToRepoNotificationMethods сериализует каналы уведомлений в JSON для колонки jsonb.
func ToRepoNotificationMethods(methods []model.NotificationMethod) ([]byte, error) {
	repoMethods := make([]repoModel.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		repoMethods = append(repoMethods, repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return json.Marshal(repoMethods)
}

// ToModelUser конвертирует модель репозитория пользователя в доменную модель.
func ToModelUser(user repoModel.User) (model.User, error) {
	userUUID, err := uuid.Parse(user.UUID)
	if err != nil {
		return model.User{}, fmt.Errorf("некорректный uuid пользователя %q: %w", user.UUID, err)
	}

	var repoMethods []repoModel.NotificationMethod
	if err := json.Unmarshal(user.NotificationMethods, &repoMethods); err != nil {
		return model.User{}, fmt.Errorf("некорректные каналы уведомлений: %w", err)
	}

	methods := make([]model.NotificationMethod, 0, len(repoMethods))
	for _, method := range repoMethods {
		methods = append(methods, model.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return model.User{
		UUID: userUUID,
		Info: model.UserInfo{
			Login:               user.Login,
			Email:               user.Email,
			NotificationMethods: methods,
		},
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}
