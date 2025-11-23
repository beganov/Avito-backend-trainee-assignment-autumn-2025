package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/api"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/database"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err, "No .env file found")
		// logger.Fatal(err, "No .env file found")
	}
	metrics.Init()
	config.VarsInit()
	handleSignals(cancel)

	database.RunMigrations(config.PostgresURL)

	database.InitDB(ctx, config.PostgresURL)
	defer database.DB.Close()

	handler := api.NewHandler(ctx)

	cache.InitCache()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/team/add", handler.AddTeam)
	e.GET("/team/get", handler.GetTeam)

	e.POST("/users/setIsActive", handler.SetUserIsActive)
	e.GET("/users/getReview", handler.GetUserReview)

	e.POST("/pullRequest/create", handler.CreatePullRequest)
	e.POST("/pullRequest/merge", handler.MergePullRequest)
	e.POST("/pullRequest/reassign", handler.ReassignPullRequest)

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down services")

	gracefulShutdown(e, database.DB)
	log.Print("App stopped")
}

func handleSignals(cancel context.CancelFunc) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		log.Print("Caught signal", sig)
		//log.Info("Caught signal", sig)
		cancel()
	}()
}

func gracefulShutdown(e *echo.Echo, db *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Print(err, "HTTP server shutdown failed")
		//log.Error(err, "HTTP server shutdown failed")
	}
	db.Close()
	log.Print("All services stopped")
}
