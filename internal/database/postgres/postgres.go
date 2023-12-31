package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(DSN string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), DSN)
	if err != nil {
		log.Fatal("Error occured while established connection to database", err)
	}

	connect, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the db pool")
	}
	defer connect.Release()

	err = connect.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return pool, err
}
