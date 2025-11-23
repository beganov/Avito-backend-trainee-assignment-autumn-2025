package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

func TestFullFlow(t *testing.T) {

	t.Run("HealthCheck", func(t *testing.T) {

		resp, err := http.Get("http://localhost:8080/health")

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

	})

	t.Run("CreateTeam", func(t *testing.T) {

		team := models.Team{

			TeamName: "dev_team",

			Members: []models.TeamMember{

				{UserID: "dev1", Username: "Developer1", IsActive: true},

				{UserID: "dev2", Username: "Developer2", IsActive: true},

				{UserID: "dev3", Username: "Developer3", IsActive: true},

				{UserID: "dev4", Username: "Developer4", IsActive: false},
			},
		}

		body, _ := json.Marshal(team)

		resp, err := http.Post("http://localhost:8080/team/add", "application/json", bytes.NewReader(body))

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

	})

	t.Run("GetTeam", func(t *testing.T) {

		resp, err := http.Get("http://localhost:8080/team/get?team_name=dev_team")

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var team models.Team

		err = json.NewDecoder(resp.Body).Decode(&team)

		assert.NoError(t, err)

		log.Printf("Team: %+v", team)

		assert.Len(t, team.Members, 4)

	})

	t.Run("CreatePR", func(t *testing.T) {

		pr := models.PullRequestShort{

			PullRequestID: "pr_e2e_1",

			PullRequestName: "E2E Test Feature",

			AuthorID: "dev1",
		}

		body, _ := json.Marshal(pr)

		resp, err := http.Post("http://localhost:8080/pullRequest/create", "application/json", bytes.NewReader(body))

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]models.PullRequest

		json.NewDecoder(resp.Body).Decode(&response)

		assert.Len(t, response["pr"].AssignedReviewers, 2)

		assert.Equal(t, "OPEN", response["pr"].Status)

	})

	t.Run("GetUserReview", func(t *testing.T) {

		resp, err := http.Get("http://localhost:8080/users/getReview?user_id=dev2")

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "dev2", response["user_id"])

		pullRequests := response["pull_requests"].([]interface{})

		assert.Greater(t, len(pullRequests), 0)

	})

	t.Run("MergePR", func(t *testing.T) {

		mergeData := map[string]string{

			"pull_request_id": "pr_e2e_1",
		}

		body, _ := json.Marshal(mergeData)

		req, _ := http.NewRequest("POST", "http://localhost:8080/pullRequest/merge", bytes.NewReader(body))

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]models.PullRequest

		err = json.NewDecoder(resp.Body).Decode(&response)

		assert.NoError(t, err)

		assert.Equal(t, "MERGED", response["pr"].Status)

	})

	t.Run("ReassignMergedPR", func(t *testing.T) {

		reassignData := map[string]string{

			"pull_request_id": "pr_e2e_1",

			"old_reviewer_id": "dev2",
		}

		body, _ := json.Marshal(reassignData)

		req, _ := http.NewRequest("POST", "http://localhost:8080/pullRequest/reassign", bytes.NewReader(body))

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

	})

	t.Run("CreateSecondPR", func(t *testing.T) {

		pr := models.PullRequestShort{

			PullRequestID: "pr_e2e_2",

			PullRequestName: "E2E Test Feature 2",

			AuthorID: "dev1",
		}

		body, _ := json.Marshal(pr)

		resp, err := http.Post("http://localhost:8080/pullRequest/create", "application/json", bytes.NewReader(body))

		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

	})

	t.Run("ReassignReviewer", func(t *testing.T) {

		reassignData := map[string]string{

			"pull_request_id": "pr_e2e_2",

			"old_reviewer_id": "dev2",
		}

		body, _ := json.Marshal(reassignData)

		req, _ := http.NewRequest("POST", "http://localhost:8080/pullRequest/reassign", bytes.NewReader(body))

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)

		if resp.StatusCode == http.StatusOK {

			var response map[string]interface{}

			json.NewDecoder(resp.Body).Decode(&response)

			assert.NotEmpty(t, response["replaced_by"])

		} else {

			assert.Equal(t, http.StatusConflict, resp.StatusCode)

		}

	})

	t.Run("DeactivateUser", func(t *testing.T) {

		activityData := map[string]interface{}{

			"user_id": "dev3",

			"is_active": false,
		}

		body, _ := json.Marshal(activityData)

		req, _ := http.NewRequest("POST", "http://localhost:8080/users/setIsActive", bytes.NewReader(body))

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

	})

}

func TestLoad(t *testing.T) {

	const N = 300

	const workers = 10

	var wg sync.WaitGroup

	wg.Add(workers)

	successCount := 0

	var mu sync.Mutex

	teamData := map[string]interface{}{

		"team_name": "load_team",

		"members": []map[string]interface{}{

			{"user_id": "dev1", "username": "Load User 1", "is_active": true},

			{"user_id": "load2", "username": "Load User 2", "is_active": true},

			{"user_id": "load3", "username": "Load User 3", "is_active": true},
		},
	}

	body, _ := json.Marshal(teamData)

	resp, err := http.Post("http://localhost:8080/team/add", "application/json", bytes.NewReader(body))

	require.NoError(t, err)

	defer resp.Body.Close()

	if resp.StatusCode != 201 {

		t.Logf("Team might already exist, continuing...")

	}

	for i := 0; i < workers; i++ {

		go func(workerID int) {

			defer wg.Done()

			client := &http.Client{Timeout: 10 * time.Second}

			for j := 0; j < N/workers; j++ {

				prData := map[string]string{

					"pull_request_id": fmt.Sprintf("load_pr_%d_%d", workerID, j),

					"pull_request_name": fmt.Sprintf("Load Test %d-%d", workerID, j),

					"author_id": "dev1",
				}

				body, _ := json.Marshal(prData)

				resp, err := client.Post("http://localhost:8080/pullRequest/create", "application/json", bytes.NewReader(body))

				if err == nil && resp.StatusCode == 201 {

					mu.Lock()

					successCount++

					mu.Unlock()

				}

				if resp != nil {

					resp.Body.Close()

				}

				time.Sleep(200 * time.Millisecond)

			}

		}(i)

	}

	wg.Wait()

	t.Logf("Load test completed: %d/%d successful requests", successCount, N)

	require.Greater(t, successCount, 299, "Success rate should be >99.9%")

}
