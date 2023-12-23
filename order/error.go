package order

import "errors"

var (
	ErrInvalidOrderNumber  = errors.New("order number is invalid")
	ErrOrderBeenRegistered = errors.New("order has already been registered")
	ErrAnotherUserOrder    = errors.New("order belongs to another user")
	ErrCurrentUserOrder    = errors.New("order already has been registred current user")
)
