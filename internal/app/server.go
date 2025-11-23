package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/api"
)

// StartServer initializes and configures the HTTP server
func StartServer(ctx context.Context) *echo.Echo {

	handler := api.NewHandler(ctx) // Create API handler with context

	e := echo.New() // Initialize Echo framework

	// Add middleware for request logging and panic recovery
	e.Use(middleware.Logger())

	e.Use(middleware.Recover())

	// Team endpoints
	e.POST("/team/add", handler.AddTeam)

	e.GET("/team/get", handler.GetTeam)

	// Users endpoints
	e.POST("/users/setIsActive", handler.SetUserIsActive)

	e.GET("/users/getReview", handler.GetUserReview)

	// PullRequest endpoints
	e.POST("/pullRequest/create", handler.CreatePullRequest)

	e.POST("/pullRequest/merge", handler.MergePullRequest)

	e.POST("/pullRequest/reassign", handler.ReassignPullRequest)

	// System endpoints
	e.GET("/health", handler.Health)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e

}
