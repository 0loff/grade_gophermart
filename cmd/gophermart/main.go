package main

import (
	"log"

	"github.com/0loff/grade_gophermart/config"
	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/server"
)

func main() {
	cfg := config.NewConfigBuilder()

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal(err)
	}

	app := server.NewApp(cfg)

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
