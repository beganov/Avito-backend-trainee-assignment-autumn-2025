package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
	pullrequest "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/pullRequest"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{ctx: ctx}
}

func (h *Handler) AddTeam(c echo.Context) error {
	timer := prometheus.NewTimer(metrics.HttpDuration)
	defer timer.ObserveDuration()
	var bindedTeam models.Team
	var TeamResponse models.TeamResponse
	err := c.Bind(&bindedTeam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.ValidationError())
	}
	TeamResponse, err = team.Add(bindedTeam, h.ctx)
	if err != nil {
		if errors.Is(err, errs.ErrTeamExists) {
			return c.JSON(http.StatusNotFound, errs.TeamExists())
		}
		return c.JSON(http.StatusInternalServerError, errs.DatabaseError())
	}
	return c.JSON(http.StatusCreated, TeamResponse)
}

func (h *Handler) GetTeam(c echo.Context) error {
	timer := prometheus.NewTimer(metrics.HttpDuration)
	defer timer.ObserveDuration()
	team_name := c.QueryParam("team_name")
	team, err := team.Get(team_name, h.ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errs.NotFound())
		}
		return c.JSON(http.StatusInternalServerError, errs.DatabaseError())
	}
	return c.JSON(http.StatusOK, team)
}

func (h *Handler) SetUserIsActive(c echo.Context) error {
	timer := prometheus.NewTimer(metrics.HttpDuration)
	defer timer.ObserveDuration()
	var bindedUser models.UserActivity
	err := c.Bind(&bindedUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.ValidationError())
	}
	user, err := team.SetActive(bindedUser, h.ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return c.JSON(http.StatusNotFound, errs.NotFound())
		}
		return c.JSON(http.StatusInternalServerError, errs.DatabaseError())
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserReview(c echo.Context) error {
	timer := prometheus.NewTimer(metrics.HttpDuration)
	defer timer.ObserveDuration()
	user_id := c.QueryParam("user_id")
	requests := pullrequest.GetPR(h.ctx, user_id)
	return c.JSON(http.StatusOK, requests)
}

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
	return c.JSON(http.StatusCreated, request)
}

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
	return c.JSON(http.StatusCreated, request)
}
