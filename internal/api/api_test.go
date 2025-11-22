package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
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
				var response team.TeamResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "backend", response.Team.TeamName)
				assert.Len(t, response.Team.Members, 2)
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

	// Создаем тестового пользователя
	cache.UserCache.Set("u1", users.User{
		UserID:   "u1",
		Username: "Alice",
		TeamName: "backend",
		IsActive: false,
	})

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success - set user active",
			requestBody: team.UserActivity{
				UserID:   "u1",
				IsActive: true,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response team.UserResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "u1", response.User.UserID)
				assert.True(t, response.User.IsActive)
			},
		},
		{
			name: "Error - user not found",
			requestBody: team.UserActivity{
				UserID:   "non_existent",
				IsActive: true,
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name:           "Error - invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusNotFound,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
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

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestCreatePullRequest(t *testing.T) {
	e := echo.New()
	resetTestData()

	cache.TeamCache.Set("backend", team.Team{
		TeamName: "backend",
		Members: []team.TeamMember{
			{UserID: "author1", Username: "Author", IsActive: true},
			{UserID: "reviewer1", Username: "Reviewer1", IsActive: true},
			{UserID: "reviewer2", Username: "Reviewer2", IsActive: true},
		},
	})
	cache.UserCache.Set("author1", users.User{
		UserID: "author1", Username: "Author", TeamName: "backend", IsActive: true,
	})

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success - create PR",
			requestBody: pullrequest.PullRequestShort{
				PullRequestID:   "pr1",
				PullRequestName: "Test PR",
				AuthorID:        "author1",
				Status:          "OPEN",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response pullrequest.PRResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "pr1", response.PullRequest.PullRequestID)
				assert.Equal(t, "OPEN", response.PullRequest.Status)
				assert.Len(t, response.PullRequest.AssignedReviewers, 2)
			},
		},
		{
			name: "Error - PR already exists",
			requestBody: pullrequest.PullRequestShort{
				PullRequestID:   "pr1",
				PullRequestName: "Test PR",
				AuthorID:        "author1",
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
		{
			name: "Error - author not found",
			requestBody: pullrequest.PullRequestShort{
				PullRequestID:   "pr2",
				PullRequestName: "Test PR",
				AuthorID:        "unknown_author",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
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

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}

func TestMergePullRequest(t *testing.T) {
	e := echo.New()
	resetTestData()

	cache.PRcache.Set("pr1", pullrequest.PullRequest{
		PullRequestID:   "pr1",
		PullRequestName: "Test PR",
		AuthorID:        "author1",
		Status:          "OPEN",
		CreatedAt:       time.Now().UTC().Format(time.RFC3339),
	})

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Success - merge PR",
			requestBody: pullrequest.PullRequestShort{
				PullRequestID: "pr1",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var response pullrequest.PRResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "MERGED", response.PullRequest.Status)
				assert.NotEmpty(t, response.PullRequest.MergedAt)
			},
		},
		{
			name: "Error - PR not found",
			requestBody: pullrequest.PullRequestShort{
				PullRequestID: "unknown_pr",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  func(t *testing.T, rec *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/pr/merge", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := MergePullRequest(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rec)
			}
		})
	}
}
