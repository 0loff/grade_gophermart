package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/models"
	userrepo "github.com/0loff/grade_gophermart/user"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type User struct {
	ID       uint32 `db:"_id, omitempty"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	ur := &UserRepository{
		db: db,
	}

	ur.CreateTable()

	return ur
}

func (r UserRepository) CreateTable() {
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		uuid TEXT NOT NULL,
		username TEXT NOT NULL,
		hash TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL
		);`)
	if err != nil {
		logger.Log.Error("Unable to create USER table", zap.Error(err))
	}

	_, err = r.db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS username ON users (username)")
	if err != nil {
		logger.Log.Error("Unable to create unique index for username field")
	}
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	now := time.Now()
	uid := uuid.New().String()

	_, err := r.db.Exec(`INSERT INTO users(uuid, username, hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		uid, user.Username, user.Password, now.Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			logger.Log.Error("Username already exists", zap.Error(err))
			return "", userrepo.ErrUserName
		}
		logger.Log.Error("Failed to create new user", zap.Error(err))
		return "", err
	}

	return uid, err
}

func (r UserRepository) GetUser(ctx context.Context, username string) (*models.User, error) {
	var User models.User

	row := r.db.QueryRowContext(ctx, `SELECT uuid, username, hash FROM users WHERE username = $1`, username)
	if err := row.Scan(&User.ID, &User.Username, &User.Password); err != nil {
		return nil, err
	}

	return &User, nil
}
