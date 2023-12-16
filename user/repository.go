package user

import (
	"context"

	"github.com/0loff/grade_gophermart/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
}
