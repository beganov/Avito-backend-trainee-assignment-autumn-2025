package api

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"
)

// Handler manages HTTP request handlers for the PR review service
type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context) *Handler {

	return &Handler{ctx: ctx}

}

// Health проверяет работоспособность сервиса

// @Summary Health check

// @Description Проверяет работоспособность сервиса

// @Tags health

// @Produce json

// @Success 200 {string} string "OK"

// @Router /health [get]

func (h *Handler) Health(c echo.Context) error {

	timer := prometheus.NewTimer(metrics.HttpDuration) // Start timer for metrics collection

	defer timer.ObserveDuration()

	return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))

}
