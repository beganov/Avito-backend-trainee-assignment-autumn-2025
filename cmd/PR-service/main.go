package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/users"
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
	var bindedTeam team.Team
	err := c.Bind(&bindedTeam)
	if err != nil { //ошибка валидации не предусмотрена - надо подумать
		return c.JSON(http.StatusBadRequest, errs.TeamExists())
	}
	err = team.Add(bindedTeam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.TeamExists())
	}
	return c.JSON(http.StatusCreated, bindedTeam)
}

func getTeam(c echo.Context) error {
	team_name := c.QueryParam("team_name")
	team, err := team.Get(team_name)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusOK, team)
}

func setUserIsActive(c echo.Context) error {
	var bindedUser users.UserActivity
	err := c.Bind(&bindedUser)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	user, err := users.SetActive(bindedUser)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusOK, user)
}

func getUserReview(c echo.Context) error {
	user_id := c.QueryParam("user_id")
	requests := users.GetPR(user_id)
	return c.JSON(http.StatusOK, requests)
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
