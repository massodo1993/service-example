package model

import (
	"time"

	"github.com/google/uuid"
)

// Session — сессия пользователя, хранится в Redis.
type Session struct {
	UUID      uuid.UUID
	UserUUID  uuid.UUID
	ExpiresAt time.Time
}
