package order

import (
	"context"

	"github.com/0loff/grade_gophermart/models"
)

type Repository interface {
	InsertOrder(ctx context.Context, order models.Order) error
	GetUUIDByOrder(ctx context.Context, order string) (string, error)
	GetOrdersByUUID(ctx context.Context, uuid string) ([]models.Order, error)
	GetBalance(ctx context.Context, uuid string) (models.Balance, error)
	GetDrawalsByUUID(ctx context.Context, uuid string) ([]models.Drawall, error)
}
