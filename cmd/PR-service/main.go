package main

import (
	"errors"
	"log/slog"
	"net/http"

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
	e.POST("/team/add", addTeam)
	e.GET("/team/get", getTeam)

	e.POST("/users/setIsActive", setUserIsActive)
	e.GET("/users/getReview", getUserReview)

	e.POST("/pullRequest/create", createPullRequest)
	e.POST("/pullRequest/merge", mergePullRequest)
	e.POST("/pullRequest/reassign", reassignPullRequest)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

// Handler
func addTeam(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func getTeam(c echo.Context) error {
	team_name := c.QueryParam("team_name")

	return c.String(http.StatusOK, "Hello, World!"+team_name)
}

func setUserIsActive(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func getUserReview(c echo.Context) error {
	user_id := c.QueryParam("user_id")
	return c.String(http.StatusOK, "Hello, World!"+user_id)
}

func createPullRequest(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func mergePullRequest(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func reassignPullRequest(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
