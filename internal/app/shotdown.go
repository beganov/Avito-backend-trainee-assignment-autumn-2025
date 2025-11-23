package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
)

// HandleSignals sets up signal handling for graceful shutdown
func HandleSignals(cancel context.CancelFunc) {

	go func() {

		c := make(chan os.Signal, 1) // Create channel to receive OS signals

		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM) // Register for interrupt and terminate signals

		sig := <-c // Wait for signal

		logger.Info("Caught signal", sig)

		cancel()

	}()

}

// GracefulShutdown performs orderly shutdown of all services
func GracefulShutdown(e *echo.Echo, db *pgxpool.Pool) {

	logger.Info("Shutting down services")

	// Create shutdown context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Shutdown HTTP server
	if err := e.Shutdown(ctx); err != nil {

		logger.Error(err, "HTTP server shutdown failed")

	}

	// Close database connection pool
	db.Close()

	logger.Info("All services stopped")

}
