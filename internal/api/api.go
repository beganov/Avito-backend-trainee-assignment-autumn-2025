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

// AddTeam создает новую команду с участниками
// @Summary Создать команду с участниками (создаёт/обновляет пользователей)
// @Tags Teams
// @Accept json
// @Produce json
// @Param team body models.Team true "Данные команды"
// @Success 201 {object} object{team=models.Team} "Команда создана"
// @Failure 400 {object} errs.ErrorResponse "Команда уже существует"
// @Router /team/add [post]
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
			return c.JSON(http.StatusBadRequest, errs.TeamExists())
		}
		return c.JSON(http.StatusInternalServerError, errs.DatabaseError())
	}
	return c.JSON(http.StatusCreated, TeamResponse)
}

// GetTeam получает информацию о команде с участниками
// @Summary Получить команду с участниками
// @Tags Teams
// @Produce json
// @Param team_name query string true "Уникальное имя команды"
// @Success 200 {object} models.Team "Объект команды"
// @Failure 404 {object} errs.ErrorResponse "Команда не найдена"
// @Router /team/get [get]
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

// SetUserIsActive обновляет статус активности пользователя
// @Summary Установить флаг активности пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param user body object{user_id=string,is_active=bool} true "Данные активности пользователя"
// @Success 200 {object} object{user=models.User} "Обновлённый пользователь"
// @Failure 404 {object} errs.ErrorResponse "Пользователь не найден"
// @Router /users/setIsActive [post]
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

// GetUserReview получает пул-реквесты для ревью пользователя
// @Summary Получить PR'ы, где пользователь назначен ревьювером
// @Tags Users
// @Produce json
// @Param user_id query string true "Идентификатор пользователя"
// @Success 200 {object} object{user_id=string,pull_requests=[]models.PullRequestShort} "Список PR'ов пользователя"
// @Router /users/getReview [get]
func (h *Handler) GetUserReview(c echo.Context) error {
	timer := prometheus.NewTimer(metrics.HttpDuration)
	defer timer.ObserveDuration()
	user_id := c.QueryParam("user_id")
	requests := pullrequest.GetPR(h.ctx, user_id)
	return c.JSON(http.StatusOK, requests)
}

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
	return c.JSON(http.StatusCreated, request)
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

// Health проверяет работоспособность сервиса
// @Summary Health check
// @Description Проверяет работоспособность сервиса
// @Tags health
// @Produce json
// @Success 200 {string} string "OK"
// @Router /health [get]
func (h *Handler) Health(c echo.Context) error {
	timer := prometheus.NewTimer(metrics.HttpDuration)
	defer timer.ObserveDuration()
	return c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
}
