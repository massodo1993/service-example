package model

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this login or email already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrSessionNotFound    = errors.New("session not found or expired")
	ErrInvalidArgument    = errors.New("login, email and password must not be empty")
)
