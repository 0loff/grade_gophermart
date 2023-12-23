package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserName           = errors.New("username already exists")
	ErrInvalidPassword    = errors.New("incorrect password")
	ErrInvalidAccessToken = errors.New("invalid access token")
)
