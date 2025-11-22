package pullrequest

import (
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/users"
)

var PRcache map[string]PullRequest = make(map[string]PullRequest)
var UserPRcache map[string][]PullRequest = make(map[string][]PullRequest)

type PullRequest struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"` // OPEN, MERGED
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         string   `json:"createdAt,omitempty"`
	MergedAt          string   `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"` // OPEN, MERGED
}

type PRResponse struct {
	PullRequest PullRequest `json:"pr"`
}

func Create(bindedPR PullRequestShort) (PRResponse, error) {
	req, ok := PRcache[bindedPR.PullRequestID]
	if ok {
		return PRResponse{}, errs.ErrPRExists
	}
	author, ok := users.UserCache[bindedPR.AuthorID]
	if !ok {
		return PRResponse{}, errs.ErrNotFound
	}
	req = PullRequest{
		PullRequestID:     bindedPR.PullRequestID,
		PullRequestName:   bindedPR.PullRequestName,
		AuthorID:          bindedPR.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: []string{},
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
	}

	reqTeam := team.TeamCache[author.TeamName]
	counter := 0
	for _, j := range reqTeam.Members {
		if j.UserID == author.UserID {
			continue
		}
		if j.IsActive {
			counter++
			req.AssignedReviewers = append(req.AssignedReviewers, j.UserID)
			UserPRcache[j.UserID] = append(UserPRcache[j.UserID], req)
		}
		if counter == 2 {
			break
		}
	}
	PRcache[bindedPR.PullRequestID] = req
	return PRResponse{PullRequest: req}, nil
}

func Merge(bindedPR PullRequestShort) (PRResponse, error) {
	req, ok := PRcache[bindedPR.PullRequestID]
	if !ok {
		return PRResponse{}, errs.ErrNotFound
	}
	if req.Status == "MERGED" {
		return PRResponse{PullRequest: req}, nil
	}
	req.Status = "MERGED"
	req.MergedAt = time.Now().UTC().Format(time.RFC3339)
	return PRResponse{PullRequest: req}, nil
}

func Reassign(bindedPR PullRequestShort) (PRResponse, error) {
	req, ok := PRcache[bindedPR.PullRequestID]
	if !ok {
		return PRResponse{}, errs.ErrNotFound
	}
	reviewer, ok := users.UserCache[bindedPR.AuthorID]
	if !ok {
		return PRResponse{}, errs.ErrNotFound
	}
	if req.Status == "MERGED" {
		return PRResponse{}, errs.ErrPRMerged
	}
	author := req.AuthorID
	if len(req.AssignedReviewers) == 0 {
		return PRResponse{}, errs.ErrNotAssigned
	}
	index := 0
	rew0 := req.AssignedReviewers[0]
	if len(req.AssignedReviewers) == 1 && rew0 != reviewer.UserID {
		return PRResponse{}, errs.ErrNotAssigned
	}
	if rew0 != reviewer.UserID {
		index++
	}
	rew1 := ""
	if len(req.AssignedReviewers) == 2 {
		rew1 = req.AssignedReviewers[1]
		if rew1 != reviewer.UserID {
			return PRResponse{}, errs.ErrNotAssigned
		}
	}
	reqTeam := team.TeamCache[reviewer.TeamName]
	for _, j := range reqTeam.Members {
		if j.UserID == author || j.UserID == rew0 || j.UserID == rew1 {
			continue
		}
		if j.IsActive {
			req.AssignedReviewers[index] = j.UserID
			for i, j := range UserPRcache[reviewer.UserID] {
				if j.PullRequestID == req.PullRequestID {
					UserPRcache[reviewer.UserID][i] = UserPRcache[reviewer.UserID][len(UserPRcache[reviewer.UserID])-1]
					UserPRcache[reviewer.UserID] = UserPRcache[reviewer.UserID][:len(UserPRcache[reviewer.UserID])-1]
					break
				}
			}
			UserPRcache[j.UserID] = append(UserPRcache[j.UserID], req)
			return PRResponse{PullRequest: req}, nil
		}
	}
	return PRResponse{}, errs.ErrNoCandidate
}

func GetPR(UserID string) UserRequests {
	return UserRequests{UserID: UserID, Pull_requests: UserPRcache[UserID]}
}

type UserRequests struct {
	UserID        string        `json:"user_id"`
	Pull_requests []PullRequest `json:"pull_requests"`
}
