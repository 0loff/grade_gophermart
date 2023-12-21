package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/0loff/grade_gophermart/balance"
	balancehttp "github.com/0loff/grade_gophermart/balance/delivery/http"
	balanceusecase "github.com/0loff/grade_gophermart/balance/usecase"
	"github.com/0loff/grade_gophermart/config"
	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/order"

	orderhttp "github.com/0loff/grade_gophermart/order/delivery/http"
	orderpostgres "github.com/0loff/grade_gophermart/order/repository/postgres"
	orderusecase "github.com/0loff/grade_gophermart/order/usecase"
	"github.com/0loff/grade_gophermart/user"
	userhttp "github.com/0loff/grade_gophermart/user/delivery/http"
	userpostgres "github.com/0loff/grade_gophermart/user/repository/postgres"
	userusecase "github.com/0loff/grade_gophermart/user/usecase"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	httpServer *http.Server

	userUC    user.UseCase
	orderUC   order.UseCase
	balanceUC balance.UseCase
}

func NewApp(cfg config.Config) *App {
	dbpool, err := initDB(cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Error("Unable to create database instance", zap.Error(err))
	}

	userRepo := userpostgres.NewUserRepository(dbpool)
	orderRepo := orderpostgres.NewOrderRepository(dbpool)

	return &App{
		userUC: userusecase.NewUserUseCase(
			userRepo,
			[]byte(cfg.SigningKey),
			time.Hour*3,
		),
		orderUC: orderusecase.NewOrderUseCase(
			orderRepo,
		),
		balanceUC: balanceusecase.NewBalanceUseCase(
			orderRepo,
		),
	}
}

func (a *App) Run(cfg config.Config) error {
	router := chi.NewRouter()
	router.Use(logger.RequestLogger)

	userhttp.RegisterHTTPEndpoints(router, a.userUC)
	authMiddleware := userhttp.NewAuthMiddleware(a.userUC).Handle

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		orderhttp.RegisterHTTPEndpoints(r, a.orderUC)
		balancehttp.RegisterHTTPEndpoints(r, a.balanceUC)
	})

	a.httpServer = &http.Server{
		Addr:    cfg.Host,
		Handler: router,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	logger.Sugar.Infoln("Host", cfg.Host)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB(DSN string) (*pgxpool.Pool, error) {
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
