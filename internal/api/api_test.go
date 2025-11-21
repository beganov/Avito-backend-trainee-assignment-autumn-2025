package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	pullrequest "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/pullRequest"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/users"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func resetTestData() {
	// Здесь будем сбрасывать тестовые данные
	// Пока что оставим пустым - реализуем позже
}

func TestAddTeam(t *testing.T) {
	e := echo.New()
	resetTestData()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success - create new team",
			requestBody: team.Team{
				TeamName: "backend",
				Members: []team.TeamMember{
					{UserID: "u1", Username: "Alice", IsActive: true},
					{UserID: "u2", Username: "Bob", IsActive: true},
				},
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response team.Team
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "backend", response.TeamName)
				assert.Len(t, response.Members, 2)
			},
		},
		{
			name: "Error - team already exists",
			requestBody: team.Team{
				TeamName: "backend", // Та же команда
				Members:  []team.TeamMember{},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				// Должны получить ошибку о существующей команде
			},
		},
		{
			name:           "Error - invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := AddTeam(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.checkResponse(t, rec)
		})
	}
}

func TestGetTeam(t *testing.T) {
	e := echo.New()
	resetTestData()

	// Сначала создаем команду для теста
	team.Add(team.Team{
		TeamName: "frontend",
		Members: []team.TeamMember{
			{UserID: "u3", Username: "Charlie", IsActive: true},
		},
	})

	tests := []struct {
		name           string
		teamName       string
		expectedStatus int
	}{
		{
			name:           "Success - get existing team",
			teamName:       "frontend",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - team not found",
			teamName:       "nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Error - empty team name",
			teamName:       "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+tt.teamName, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := GetTeam(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestSetUserIsActive(t *testing.T) {
	e := echo.New()
	resetTestData()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "Success - set user active",
			requestBody: users.UserActivity{
				UserID:   "u1",
				IsActive: true,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - invalid user data",
			requestBody:    "invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/user/active", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := SetUserIsActive(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestCreatePullRequest(t *testing.T) {
	e := echo.New()
	resetTestData()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "Success - create PR",
			requestBody: pullrequest.PullRequestShort{
				PullRequestID:   "pr1",
				PullRequestName: "Fix bug",
				AuthorID:        "u1",
				Status:          "open",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Error - invalid PR data",
			requestBody:    "invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/pr/create", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := CreatePullRequest(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
