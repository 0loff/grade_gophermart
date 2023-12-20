package user

import "context"

type UseCase interface {
	SignUp(ctx context.Context, username, password string) (string, error)
	SignIn(ctx context.Context, username, password string) (string, error)
}
