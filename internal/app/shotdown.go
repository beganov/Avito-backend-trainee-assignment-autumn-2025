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

func HandleSignals(cancel context.CancelFunc) {

	go func() {

		c := make(chan os.Signal, 1)

		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		sig := <-c

		logger.Info("Caught signal", sig)

		cancel()

	}()

}

func GracefulShutdown(e *echo.Echo, db *pgxpool.Pool) {

	logger.Info("Shutting down services")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := e.Shutdown(ctx); err != nil {

		logger.Error(err, "HTTP server shutdown failed")

	}

	db.Close()

	logger.Info("All services stopped")

}
