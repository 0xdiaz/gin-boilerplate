package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database/migrations"
	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database/seeders"
	"github.com/0xdiaz/gin-boilerplate/internal/app/routers"
	"github.com/0xdiaz/gin-boilerplate/pkg/config"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/0xdiaz/gin-boilerplate/pkg/metrics"
)

func main() {
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	cfg := config.Get()
	if cfg == nil {
		logger.Fatalf("config not loaded")
	}
	loc, err := time.LoadLocation(cfg.Server.Timezone)
	if err != nil {
		logger.Fatalf("invalid SERVER_TIMEZONE %q: %v", cfg.Server.Timezone, err)
	}
	time.Local = loc

	// Initialize metrics tracking
	metrics.Init()
	logger.Infof("Metrics tracking initialized")

	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	if err := migrations.Migrate(); err != nil {
		logger.Fatalf("database migration failed: %v", err)
	}

	if config.IsDevelopment() {
		if err := seeders.Run(); err != nil {
			logger.Fatalf("database seeding failed: %v", err)
		}
	}

	router := routers.SetupRoute()

	addr := config.ServerConfig()
	requestTimeout := time.Duration(cfg.Server.RequestTimeoutSeconds) * time.Second
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  requestTimeout,
		WriteTimeout: requestTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server listen error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Infof("shutdown signal received, shutting down gracefully")

	shutdownTimeout := time.Duration(cfg.Server.ShutdownTimeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown error: %v", err)
	}

	if sqlDB, err := database.GetDB().DB(); err == nil && sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			logger.Errorf("database close error: %v", err)
		}
	}

	logger.Infof("shutdown complete")
}
