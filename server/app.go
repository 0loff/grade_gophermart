package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/0loff/grade_gophermart/config"
	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/user"
	userhttp "github.com/0loff/grade_gophermart/user/delivery/http"
	userpostgres "github.com/0loff/grade_gophermart/user/repository/postgres"
	userusecase "github.com/0loff/grade_gophermart/user/usecase"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	httpServer *http.Server

	userUC user.UseCase
}

func NewApp(cfg config.Config) *App {
	db, err := initDB(cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Error("Unable to create database instance", zap.Error(err))
	}

	userRepo := userpostgres.NewUserRepository(db)

	return &App{
		userUC: userusecase.NewUserUseCase(
			userRepo,
			[]byte(cfg.SigningKey),
			time.Hour*3,
		),
	}
}

func (a *App) Run(cfg config.Config) error {
	router := chi.NewRouter()
	router.Use(logger.RequestLogger)

	userhttp.RegisterHTTPEndpoints(router, a.userUC)

	router.Route("/api/user/", func(router chi.Router) {

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

func initDB(DSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		log.Fatal("Error occured while established connection to database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}
