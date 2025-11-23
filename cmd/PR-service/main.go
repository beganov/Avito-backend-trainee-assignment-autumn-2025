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
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	LoadCacheFromDB(ctx)

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

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

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

func LoadCacheFromDB(ctx context.Context) error {
	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)
	defer cancel()

	rows, err := database.DB.Query(dbCtx, `SELECT team_name FROM teams`)
	if err != nil {
		log.Printf("failed to load teams for cache: %v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var teamName string
		if err := rows.Scan(&teamName); err != nil {
			log.Printf("failed to scan team name: %v", err)
			continue
		}

		team, err, _ := database.GetTeamFromDB(dbCtx, teamName)
		if err != nil {
			log.Printf("failed to load team %s: %v", teamName, err)
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
