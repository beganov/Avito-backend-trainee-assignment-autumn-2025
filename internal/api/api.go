package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
	pullrequest "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/pullRequest"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ctx context.Context
}

func NewHandler(ctx context.Context) *Handler {
	return &Handler{ctx: ctx}
}

func (h *Handler) AddTeam(c echo.Context) error {
	var bindedTeam models.Team
	var TeamResponse models.TeamResponse
	err := c.Bind(&bindedTeam)
	if err != nil { //ошибка валидации не предусмотрена - надо подумать
		return c.JSON(http.StatusBadRequest, errs.NotFound())
	}
	TeamResponse, err = team.Add(bindedTeam, h.ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.ErrTeamExists)
	}
	return c.JSON(http.StatusCreated, TeamResponse)
}

func (h *Handler) GetTeam(c echo.Context) error {
	team_name := c.QueryParam("team_name")
	team, err := team.Get(team_name, h.ctx)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusOK, team)
}

func (h *Handler) SetUserIsActive(c echo.Context) error {
	var bindedUser models.UserActivity
	err := c.Bind(&bindedUser)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	user, err := team.SetActive(bindedUser, h.ctx)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserReview(c echo.Context) error {
	user_id := c.QueryParam("user_id")
	requests := pullrequest.GetPR(h.ctx, user_id)
	return c.JSON(http.StatusOK, requests)
}

func (h *Handler) CreatePullRequest(c echo.Context) error {
	var bindedPR models.PullRequestShort
	err := c.Bind(&bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	request, err := pullrequest.Create(h.ctx, bindedPR)
	if err != nil {
		if errors.Is(err, errs.ErrPRExists) {
			return c.JSON(http.StatusConflict, errs.PRExists())
		}
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusCreated, request)
}

func (h *Handler) MergePullRequest(c echo.Context) error {
	var bindedPR models.PullRequestShort
	err := c.Bind(&bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	request, err := pullrequest.Merge(h.ctx, bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusCreated, request)
}

func (h *Handler) ReassignPullRequest(c echo.Context) error {
	var bindedPR models.PRReassign
	err := c.Bind(&bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
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
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusCreated, request)
}
