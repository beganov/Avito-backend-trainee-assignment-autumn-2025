package api

import (
	"errors"
	"net/http"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	pullrequest "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/pullRequest"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
	"github.com/labstack/echo/v4"
)

func AddTeam(c echo.Context) error {
	var bindedTeam team.Team
	var TeamResponse team.TeamResponse
	err := c.Bind(&bindedTeam)
	if err != nil { //ошибка валидации не предусмотрена - надо подумать
		return c.JSON(http.StatusBadRequest, errs.NotFound())
	}
	TeamResponse, err = team.Add(bindedTeam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errs.TeamExists())
	}
	return c.JSON(http.StatusCreated, TeamResponse)
}

func GetTeam(c echo.Context) error {
	team_name := c.QueryParam("team_name")
	team, err := team.Get(team_name)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusOK, team)
}

func SetUserIsActive(c echo.Context) error {
	var bindedUser team.UserActivity
	err := c.Bind(&bindedUser)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	user, err := team.SetActive(bindedUser)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusOK, user)
}

func GetUserReview(c echo.Context) error {
	user_id := c.QueryParam("user_id")
	requests := pullrequest.GetPR(user_id)
	return c.JSON(http.StatusOK, requests)
}

func CreatePullRequest(c echo.Context) error {
	var bindedPR pullrequest.PullRequestShort
	err := c.Bind(&bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	request, err := pullrequest.Create(bindedPR)
	if err != nil {
		if errors.Is(err, errs.ErrPRExists) {
			return c.JSON(http.StatusConflict, errs.PRExists())
		}
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusCreated, request)
}

func MergePullRequest(c echo.Context) error {
	var bindedPR pullrequest.PullRequestShort
	err := c.Bind(&bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	request, err := pullrequest.Merge(bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	return c.JSON(http.StatusCreated, request)
}

func ReassignPullRequest(c echo.Context) error {
	var bindedPR pullrequest.PRReassign
	err := c.Bind(&bindedPR)
	if err != nil {
		return c.JSON(http.StatusNotFound, errs.NotFound())
	}
	request, err := pullrequest.Reassign(bindedPR)
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
