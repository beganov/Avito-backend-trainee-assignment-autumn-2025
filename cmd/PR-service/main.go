// @title PR Reviewer Assignment Service (Test Task, Fall 2025)

// @version 1.0

// @description Сервис назначения ревьюеров для Pull Request’ов

// @contact.name API Support

// @contact.url http://www.swagger.io/support

// @contact.email support@swagger.io

// @license.name Apache 2.0

// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

// @BasePath /

// @schemes http

// @produce json

// @consumes json

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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/docs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/api"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/database"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	logger.Init("info")

	if err := godotenv.Load(); err != nil {

		logger.Fatal(err, "No .env file found")

	}

	metrics.Init()

	config.VarsInit()

	handleSignals(cancel)

	database.RunMigrations(config.PostgresURL)

	database.InitDB(ctx, config.PostgresURL)

	defer database.DB.Close()

	handler := api.NewHandler(ctx)

	cache.InitCache()

	err := LoadCacheFromDB(ctx)
	if err != nil {
		logger.Error(err, "cache dont loaded")

	}

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

	e.GET("/health", handler.Health)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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

		logger.Info("Caught signal", sig)

		cancel()

	}()

}

func gracefulShutdown(e *echo.Echo, db *pgxpool.Pool) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := e.Shutdown(ctx); err != nil {

		logger.Error(err, "HTTP server shutdown failed")

	}

	db.Close()

	log.Print("All services stopped")

}

func LoadCacheFromDB(ctx context.Context) error {

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)

	defer cancel()

	rows, err := database.DB.Query(dbCtx, `SELECT team_name FROM teams`)

	if err != nil {

		logger.Info("failed to load teams for cache: %v", err)

		return err

	}

	defer rows.Close()

	for rows.Next() {

		var teamName string

		if err := rows.Scan(&teamName); err != nil {

			logger.Info("failed to scan team name: %v", err)

			continue

		}

		team, err, _ := database.GetTeamFromDB(dbCtx, teamName)

		if err != nil {

			logger.Info("failed to load team %s: %v", teamName, err)

			continue

		}

		cache.TeamCache.Set(teamName, team)

	}

	userRows, err := database.DB.Query(dbCtx, `SELECT user_id FROM users`)

	if err != nil {

		log.Printf("failed to load users for cache: %v", err)

		return err

	}

	defer userRows.Close()

	for userRows.Next() {

		var userID string

		if err := userRows.Scan(&userID); err != nil {

			log.Printf("failed to scan user id: %v", err)

			continue

		}

		user, err, _ := database.GetUserFromDB(dbCtx, userID)

		if err != nil {

			log.Printf("failed to load user %s: %v", userID, err)

			continue

		}

		cache.UserCache.Set(userID, user)

	}

	prRows, err := database.DB.Query(dbCtx, `SELECT pull_request_id FROM pull_requests`)

	if err != nil {

		log.Printf("failed to load PRs for cache: %v", err)

		return err

	}

	defer prRows.Close()

	for prRows.Next() {

		var prID string

		if err := prRows.Scan(&prID); err != nil {

			log.Printf("failed to scan PR id: %v", err)

			continue

		}

		pr, err, _ := database.GetPRFromDB(dbCtx, prID)

		if err != nil {

			log.Printf("failed to load PR %s: %v", prID, err)

			continue

		}

		cache.PRcache.Set(prID, pr)

	}

	return nil

}
