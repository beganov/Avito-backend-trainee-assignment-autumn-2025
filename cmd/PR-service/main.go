package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/team/add", api.AddTeam)
	e.GET("/team/get", api.GetTeam)

	e.POST("/users/setIsActive", api.SetUserIsActive)
	e.GET("/users/getReview", api.GetUserReview)

	e.POST("/pullRequest/create", api.CreatePullRequest)
	e.POST("/pullRequest/merge", api.MergePullRequest)
	e.POST("/pullRequest/reassign", api.ReassignPullRequest)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
