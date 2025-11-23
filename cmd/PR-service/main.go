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
	"net/http"

	_ "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/docs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/app"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/database"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	app.Init() // Initialize  components (dotenv, config, metrics)

	database.RunMigrations(config.PostgresURL) // Run database migrations

	database.InitDB(ctx, config.PostgresURL) // Initialize database connection

	defer database.DB.Close()

	err := app.LoadCacheFromDB(ctx) // Load data from database into cache

	if err != nil {
		logger.Error(err, "cache dont loaded") // if cache load error - work continue
	}

	e := app.StartServer(ctx) // Setup and configure HTTP server

	go func() {

		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {

			logger.Error(err, "failed to start server")

		}

	}()

	<-ctx.Done() // Wait for shutdown signal (SIGINT, SIGTERM)

	app.GracefulShutdown(e, database.DB) // Graceful shutdown

}
