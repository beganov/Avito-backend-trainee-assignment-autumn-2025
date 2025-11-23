package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
)

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
