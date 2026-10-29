package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

type User struct {
	UUID         uuid.UUID
	Info         UserInfo
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRegistrationInfo — данные для регистрации пользователя.
type UserRegistrationInfo struct {
	Info     UserInfo
	Password string
}
