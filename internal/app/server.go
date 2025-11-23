package app

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/api"
)

func StartServer(ctx context.Context) *echo.Echo {

	handler := api.NewHandler(ctx)

	e := echo.New()

	e.Use(middleware.Logger())

	e.Use(middleware.Recover())

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

	return e

}
