package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vigorouzis/aibolit-notification/config"
	"github.com/vigorouzis/aibolit-notification/internal/core/application"
	"github.com/vigorouzis/aibolit-notification/internal/infrastructure/postgres"
	"github.com/vigorouzis/aibolit-notification/internal/interface/http"
	"log/slog"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.Default()

	cfg, err := config.FromENV()
	if err != nil {
		logger.Error("Failed to load configuration from environment", "error", err)
		os.Exit(1)
	}
	logger.Info("Configuration loaded successfully", "config", cfg)

	db, err := postgres.ConnectViaConfig(cfg.Postgres)
	if err != nil {
		logger.Error("Failed to connect to the Postgres database", "error", err, "PostgresConfig", cfg.Postgres)
		os.Exit(1)
	}
	logger.Info("Connected to the Postgres database", "PostgresConfig", cfg.Postgres)

	if err := db.PingContext(ctx); err != nil {
		logger.Error("Database did not respond to ping", "error", err)
		os.Exit(1)
	}
	logger.Info("Database ping successful")

	cl := postgres.New(db)
	logger.Info("Postgres repository initialized")

	service := application.NewService(cl)
	logger.Info("Schedule service initialized")

	server, err := http.New(cfg.HTTP, service, logger)
	if err != nil {
		logger.Error("Failed to initialize the HTTP server", "error", err, "HTTPConfig", cfg.HTTP)
		os.Exit(1)
	}
	logger.Info("HTTP server initialized", "HTTPConfig", cfg.HTTP)

	logger.Info("Starting the HTTP server...")
	if err = server.Run(ctx); err != nil {
		logger.Error("HTTP server encountered an error during runtime", "error", err)
		os.Exit(1)
	}
	logger.Info("HTTP server stopped gracefully")
}
