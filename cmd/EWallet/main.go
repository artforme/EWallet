package main

import (
	"EWallet/internal/config"
	"EWallet/internal/http-server/handlers/createWallet"
	"EWallet/internal/http-server/handlers/showHistory"
	"EWallet/internal/http-server/handlers/showWallet"
	"EWallet/internal/http-server/handlers/transfer"
	mwLogger "EWallet/internal/http-server/middleware/logger"
	"EWallet/internal/loggers"
	"EWallet/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
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
	storage, err := sqlite.New(os.Getenv("STORAGE_PATH"))
	if err != nil {
		log.Error("failed to create storage", slog.StringValue(err.Error()))
		os.Exit(1)
	}

	router := chi.NewRouter()
	// using middleware RequestID to know ID client in log
	router.Use(middleware.RequestID)
	// using our customized log
	router.Use(mwLogger.New(log))
	// this middleware recover process if some handler caught a mistake
	router.Use(middleware.Recoverer)
	// this middleware allows us read some info form url
	router.Use(middleware.URLFormat)

	// There are some handlers that we use in specific method and url
	router.Post("/api/v1/wallet", createWallet.New(log, storage))
	router.Post("/api/v1/wallet/{walletId}/send", transfer.New(log, storage))
	router.Get("/api/v1/wallet/{walletId}/history", showHistory.New(log, storage))
	router.Get("/api/v1/wallet/{walletId}", showWallet.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	// setup server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Idle_timeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
