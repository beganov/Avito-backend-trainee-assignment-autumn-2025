package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
	pullrequest "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/pullRequest"
)

// CreatePullRequest создает новый пул-реквест

// @Summary Создать PR и автоматически назначить до 2 ревьюверов из команды автора

// @Tags PullRequests

// @Accept json

// @Produce json

// @Param pr body object{pull_request_id=string,pull_request_name=string,author_id=string} true "Данные пул-реквеста"

// @Success 201 {object} object{pr=models.PullRequest} "PR создан"

// @Failure 404 {object} errs.ErrorResponse "Автор/команда не найдены"

// @Failure 409 {object} errs.ErrorResponse "PR уже существует"

// @Router /pullRequest/create [post]

func (h *Handler) CreatePullRequest(c echo.Context) error {

	timer := prometheus.NewTimer(metrics.HttpDuration)

	defer timer.ObserveDuration()

	var bindedPR models.PullRequestShort

	err := c.Bind(&bindedPR)

	if err != nil {

		return c.JSON(http.StatusBadRequest, errs.ValidationError())

	}

	request, err := pullrequest.Create(h.ctx, bindedPR)

	if err != nil {

		if errors.Is(err, errs.ErrPRExists) {

			return c.JSON(http.StatusConflict, errs.PRExists())

		}

		if errors.Is(err, errs.ErrNotFound) {

			return c.JSON(http.StatusNotFound, errs.NotFound())

		}

		return c.JSON(http.StatusInternalServerError, errs.DatabaseError())

	}

	return c.JSON(http.StatusCreated, request)

}

// MergePullRequest мержит пул-реквест

// @Summary Пометить PR как MERGED (идемпотентная операция)

// @Tags PullRequests

// @Accept json

// @Produce json

// @Param pr body object{pull_request_id=string} true "ID пул-реквеста"

// @Success 200 {object} object{pr=models.PullRequest} "PR в состоянии MERGED"

// @Failure 404 {object} errs.ErrorResponse "PR не найден"

// @Router /pullRequest/merge [post]

func (h *Handler) MergePullRequest(c echo.Context) error {

	timer := prometheus.NewTimer(metrics.HttpDuration)

	defer timer.ObserveDuration()

	var bindedPR models.PullRequestShort

	err := c.Bind(&bindedPR)

	if err != nil {

		return c.JSON(http.StatusBadRequest, errs.ValidationError())

	}

	request, err := pullrequest.Merge(h.ctx, bindedPR)

	if err != nil {

		if errors.Is(err, errs.ErrNotFound) {

			return c.JSON(http.StatusNotFound, errs.NotFound())

		}

		return c.JSON(http.StatusInternalServerError, errs.DatabaseError())

	}

	return c.JSON(http.StatusOK, request)

}

// ReassignPullRequest переназначает пул-реквест

// @Summary Переназначить конкретного ревьювера на другого из его команды

// @Tags PullRequests

// @Accept json

// @Produce json

// @Param reassignment body object{pull_request_id=string,old_user_id=string} true "Данные для переназначения"

// @Success 200 {object} object{pr=models.PullRequest,replaced_by=string} "Переназначение выполнено"

// @Failure 404 {object} errs.ErrorResponse "PR или пользователь не найден"

// @Failure 409 {object} errs.ErrorResponse "Нарушение доменных правил переназначения"

// @Router /pullRequest/reassign [post]

func (h *Handler) ReassignPullRequest(c echo.Context) error {

	timer := prometheus.NewTimer(metrics.HttpDuration)

	defer timer.ObserveDuration()

	var bindedPR models.PRReassign

	err := c.Bind(&bindedPR)

	if err != nil {

		return c.JSON(http.StatusBadRequest, errs.ValidationError())

	}

	request, err := pullrequest.Reassign(h.ctx, bindedPR)

	if err != nil {

		if errors.Is(err, errs.ErrPRMerged) {

			return c.JSON(http.StatusConflict, errs.PRMerged())

		}

		if errors.Is(err, errs.ErrNotAssigned) {

			return c.JSON(http.StatusConflict, errs.NotAssigned())

		}

		if errors.Is(err, errs.ErrNoCandidate) {

			return c.JSON(http.StatusConflict, errs.NoCandidate())

		}

		if errors.Is(err, errs.ErrNotFound) {

			return c.JSON(http.StatusNotFound, errs.NotFound())

		}

		return c.JSON(http.StatusInternalServerError, errs.DatabaseError())

	}

	return c.JSON(http.StatusOK, request)

}
