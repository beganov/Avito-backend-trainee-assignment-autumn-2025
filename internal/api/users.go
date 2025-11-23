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
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
)

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
