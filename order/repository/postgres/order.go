package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/models"
	"github.com/0loff/grade_gophermart/order"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Order struct {
	ID   uint32 `db:"_id, omitempty"`
	Num  string `db:"order_number"`
	UUID string `db:"uuid"`
}

type OrderRepository struct {
	dbpool *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	or := &OrderRepository{
		dbpool: db,
	}

	or.CreateTable()
	return or
}

func (r OrderRepository) CreateTable() {
	_, err := r.dbpool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders (
		id serial PRIMARY KEY,
		order_num TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'NEW',
		accrual INTEGER,
		uuid text,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
		CONSTRAINT fk_users
			FOREIGN KEY(uuid)
				REFERENCES users(uuid)
				ON DELETE CASCADE
	);`)
	if err != nil {
		logger.Log.Error("Unable to create ORDERS table", zap.Error(err))
	}

	_, err = r.dbpool.Exec(context.Background(), "CREATE UNIQUE INDEX IF NOT EXISTS order_num ON orders (order_num)")
	if err != nil {
		logger.Log.Error("Unable to create unique index for order_num field")
	}
}

func (r OrderRepository) InsertOrder(ctx context.Context, orderNum, uid string) error {
	now := time.Now()

	_, err := r.dbpool.Exec(context.Background(), `INSERT INTO orders(order_num, uuid, created_at, updated_at) VALUES ($1, $2, $3, $4)`,
		orderNum, uid, now.Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code):
			logger.Log.Error("Order has already been registered", zap.Error(err))
			return order.ErrOrderBeenRegistered

		default:
			logger.Log.Error("Failed to register order.", zap.Error(err))
			return err
		}
	}

	return nil
}

func (r OrderRepository) GetUUIDByOrder(ctx context.Context, order string) (string, error) {
	row := r.dbpool.QueryRow(ctx, `SELECT uuid FROM orders WHERE order_num = $1`, order)

	var UUID string
	if err := row.Scan(&UUID); err != nil {
		logger.Log.Error("Unable to parse the received UUID by order from DB", zap.Error(err))
		return "", err
	}

	return UUID, nil
}

func (r OrderRepository) GetOrdersByUUID(ctx context.Context, uuid string) ([]models.Order, error) {
	var Orders []models.Order

	rows, err := r.dbpool.Query(ctx, `SELECT order_num, status, coalesce(accrual, 0), created_at FROM orders WHERE uuid = $1`, uuid)
	if err != nil {
		logger.Log.Error("Unrecognized data from the database \n", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			Order     models.Order
			createdAt time.Time
		)
		if err := rows.Scan(&Order.OrderNum, &Order.Status, &Order.Accrual, &createdAt); err != nil {
			logger.Log.Error("Unable to parse the received value", zap.Error(err))
			continue
		}

		Order.CreatedAt = createdAt.Format(time.RFC3339)
		Orders = append(Orders, Order)
	}

	if err = rows.Err(); err != nil {
		logger.Log.Error("Unexpected error from parse data in rows next loop", zap.Error(err))
		return Orders, err
	}

	return Orders, nil
}
