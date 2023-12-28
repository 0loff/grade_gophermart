package main

import (
	"log"

	"github.com/0loff/grade_gophermart/config"
	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/server"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewConfigBuilder()

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	app := server.NewGophermart(cfg, router)

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
