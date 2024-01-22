package main

import (
	"EWallet/internal/config"
	"EWallet/internal/loggers"
	"EWallet/internal/storage/sqlite"
	"log/slog"
	"os"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup logger depending on env
	log := loggers.SetupLogger(cfg.Env)

	//check logger on different levels
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	//init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to create storage", slog.StringValue(err.Error()))
		os.Exit(1)
	}

	_ = storage
}
