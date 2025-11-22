package pullrequest

import (
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/team"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/users"
)

var UserPRcache map[string][]PullRequestShort = make(map[string][]PullRequestShort)

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
	_, ok := cache.PRcache.Get(bindedPR.PullRequestID)
	if ok {
		return PRResponse{}, errs.ErrPRExists
	}
	iUser, ok := cache.UserCache.Get(bindedPR.AuthorID)
	if !ok {
		return PRResponse{}, errs.ErrNotFound
	}
	author := iUser.(users.User)
	req := PullRequest{
		PullRequestID:     bindedPR.PullRequestID,
		PullRequestName:   bindedPR.PullRequestName,
		AuthorID:          bindedPR.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: []string{},
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
	}
	iteam, _ := cache.TeamCache.Get(author.TeamName)
	reqTeam := iteam.(team.Team)
	counter := 0
	for _, j := range reqTeam.Members {
		if j.UserID == author.UserID {
			continue
		}
		if j.IsActive {
			counter++
			req.AssignedReviewers = append(req.AssignedReviewers, j.UserID)
			UserPRcache[j.UserID] = append(UserPRcache[j.UserID], PrToPrShort(req))
		}
		if counter == 2 {
			break
		}
	}
	cache.PRcache.Set(bindedPR.PullRequestID, req)
	return PRResponse{PullRequest: req}, nil
}

func Merge(bindedPR PullRequestShort) (PRResponse, error) {
	iPR, ok := cache.PRcache.Get(bindedPR.PullRequestID)
	if !ok {
		return PRResponse{}, errs.ErrNotFound
	}
	req := iPR.(PullRequest)
	if req.Status == "MERGED" {
		return PRResponse{PullRequest: req}, nil
	}
	req.Status = "MERGED"
	req.MergedAt = time.Now().UTC().Format(time.RFC3339)
	cache.PRcache.Set(bindedPR.PullRequestID, req)

	for _, k := range req.AssignedReviewers {
		for i, j := range UserPRcache[k] {
			if j.PullRequestID == req.PullRequestID {
				UserPRcache[k][i].Status = "MERGED"
				break
			}
		}
	}
	return PRResponse{PullRequest: req}, nil
}

func Reassign(bindedPR PRReassign) (PRReassignResponse, error) {
	iPR, ok := cache.PRcache.Get(bindedPR.PullRequestID)
	if !ok {
		return PRReassignResponse{}, errs.ErrNotFound
	}
	req := iPR.(PullRequest)
	iUser, ok := cache.UserCache.Get(bindedPR.OldReviewerID)
	reviewer := iUser.(users.User)
	if !ok {
		return PRReassignResponse{}, errs.ErrNotFound
	}

	if req.Status == "MERGED" {
		return PRReassignResponse{}, errs.ErrPRMerged
	}

	stopUserMap := make(map[string]int, 3)
	stopUserMap[req.AuthorID]++
	index := -1
	for i, j := range req.AssignedReviewers {
		stopUserMap[j]++
		if j == bindedPR.OldReviewerID {
			index = i
		}
	}

	if index == -1 {
		return PRReassignResponse{}, errs.ErrNotAssigned
	}

	iteam, _ := cache.TeamCache.Get(reviewer.TeamName)
	reqTeam := iteam.(team.Team)

	for _, k := range reqTeam.Members {
		if _, ok := stopUserMap[k.UserID]; ok {
			continue
		}
		if k.IsActive {
			req.AssignedReviewers[index] = k.UserID
			for i, j := range UserPRcache[reviewer.UserID] {
				if j.PullRequestID == req.PullRequestID {
					UserPRcache[reviewer.UserID][i] = UserPRcache[reviewer.UserID][len(UserPRcache[reviewer.UserID])-1]
					UserPRcache[reviewer.UserID] = UserPRcache[reviewer.UserID][:len(UserPRcache[reviewer.UserID])-1]
					break
				}
			}
			UserPRcache[k.UserID] = append(UserPRcache[k.UserID], PrToPrShort(req))
			return PRReassignResponse{PullRequest: req, ReplacedBy: reviewer.UserID}, nil
		}
	}
	return PRReassignResponse{}, errs.ErrNoCandidate
}

func GetPR(UserID string) UserRequests {
	return UserRequests{UserID: UserID, PullRequests: UserPRcache[UserID]}
}

func PrToPrShort(PR PullRequest) PullRequestShort {
	return PullRequestShort{
		PullRequestID:   PR.PullRequestID,
		PullRequestName: PR.PullRequestName,
		AuthorID:        PR.AuthorID,
		Status:          PR.Status,
	}
}

type UserRequests struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

type PRReassign struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type PRReassignResponse struct {
	PullRequest PullRequest `json:"pr"`
	ReplacedBy  string      `json:"replaced_by"`
}
