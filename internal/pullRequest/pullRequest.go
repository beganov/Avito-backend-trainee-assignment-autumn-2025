package pullrequest

import (
	"context"
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/database"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

// PR status constants
var MergeStatus = "MERGED"
var OpenStatus = "OPEN"

// Create creates a new pull request with automatically assigned reviewers
func Create(ctx context.Context, bindedPR models.PullRequestShort) (models.PRResponse, error) {

	_, ok := cache.PRcache.Get(bindedPR.PullRequestID)

	if ok {

		return models.PRResponse{}, errs.ErrPRExists

	}

	_, err, ok := database.GetPRFromDB(ctx, bindedPR.PullRequestID)

	if err != nil {

		return models.PRResponse{}, errs.ErrDatabase

	}

	if ok {

		return models.PRResponse{}, errs.ErrPRExists

	}

	iUser, ok := cache.UserCache.Get(bindedPR.AuthorID)

	if !ok {

		iUser, err, ok = database.GetUserFromDB(ctx, bindedPR.AuthorID)

		if err != nil {

			return models.PRResponse{}, errs.ErrDatabase

		}

		if !ok {

			return models.PRResponse{}, errs.ErrNotFound

		}

		cache.UserCache.Set(bindedPR.AuthorID, iUser.(models.User))

	}

	author := iUser.(models.User)

	req := models.PullRequest{

		PullRequestID: bindedPR.PullRequestID,

		PullRequestName: bindedPR.PullRequestName,

		AuthorID: bindedPR.AuthorID,

		Status: OpenStatus,

		AssignedReviewers: []string{},

		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	iTeam, ok := cache.TeamCache.Get(author.TeamName)

	if !ok {

		iTeam, err, ok = database.GetTeamFromDB(ctx, author.TeamName)

		if err != nil || !ok {

			return models.PRResponse{}, errs.ErrDatabase

		}

		cache.TeamCache.Set(author.TeamName, iTeam.(models.Team))

	}

	reqTeam := iTeam.(models.Team)

	counter := 0

	for _, j := range reqTeam.Members {

		if j.UserID == author.UserID {

			continue

		}

		if j.IsActive {

			counter++

			req.AssignedReviewers = append(req.AssignedReviewers, j.UserID)

		}

		if counter == 2 {

			break

		}

	}

	cache.PRcache.Set(bindedPR.PullRequestID, req)

	err = database.SetPRToDB(ctx, req)

	if err != nil {

		return models.PRResponse{}, errs.ErrDatabase

	}

	return models.PRResponse{PullRequest: req}, nil

}

// Merge updates a pull request status to MERGED (idempotent operation)
func Merge(ctx context.Context, bindedPR models.PullRequestShort) (models.PRResponse, error) {

	var err error

	iPR, ok := cache.PRcache.Get(bindedPR.PullRequestID)

	if !ok {

		iPR, err, ok = database.GetPRFromDB(ctx, bindedPR.PullRequestID)

		if err != nil {

			return models.PRResponse{}, errs.ErrDatabase

		}

		if !ok {

			return models.PRResponse{}, errs.ErrNotFound

		}

	}

	req := iPR.(models.PullRequest)

	if req.Status == MergeStatus {

		return models.PRResponse{PullRequest: req}, nil

	}

	req.Status = MergeStatus

	req.MergedAt = time.Now().UTC().Format(time.RFC3339)

	cache.PRcache.Set(bindedPR.PullRequestID, req)

	err = database.SetPRToDB(ctx, req)

	if err != nil {

		return models.PRResponse{}, errs.ErrDatabase

	}

	return models.PRResponse{PullRequest: req}, nil

}

// Reassign replaces a reviewer with another active team member
func Reassign(ctx context.Context, bindedPR models.PRReassign) (models.PRReassignResponse, error) {

	var err error

	iPR, ok := cache.PRcache.Get(bindedPR.PullRequestID)

	if !ok {

		iPR, err, ok = database.GetPRFromDB(ctx, bindedPR.PullRequestID)

		if err != nil {

			return models.PRReassignResponse{}, errs.ErrDatabase

		}

		if !ok {

			return models.PRReassignResponse{}, errs.ErrNotFound

		}

	}

	req := iPR.(models.PullRequest)

	iUser, ok := cache.UserCache.Get(bindedPR.OldReviewerID)

	if !ok {

		iUser, err, ok = database.GetUserFromDB(ctx, bindedPR.OldReviewerID)

		if err != nil {

			return models.PRReassignResponse{}, errs.ErrDatabase

		}

		if !ok {

			return models.PRReassignResponse{}, errs.ErrNotFound

		}

		cache.UserCache.Set(bindedPR.OldReviewerID, iUser.(models.User))

	}

	reviewer := iUser.(models.User)

	if req.Status == MergeStatus {

		return models.PRReassignResponse{}, errs.ErrPRMerged

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

		return models.PRReassignResponse{}, errs.ErrNotAssigned

	}

	iTeam, ok := cache.TeamCache.Get(reviewer.TeamName)

	if !ok {

		iTeam, err, ok = database.GetTeamFromDB(ctx, reviewer.TeamName)

		if err != nil || !ok {

			return models.PRReassignResponse{}, errs.ErrDatabase

		}

		cache.TeamCache.Set(reviewer.TeamName, iTeam.(models.Team))

	}

	reqTeam := iTeam.(models.Team)

	for _, k := range reqTeam.Members {

		if _, ok := stopUserMap[k.UserID]; ok {

			continue

		}

		if k.IsActive {

			req.AssignedReviewers[index] = k.UserID

			cache.PRcache.Set(bindedPR.PullRequestID, req)

			err = database.SetPRToDB(ctx, req)

			if err != nil {

				return models.PRReassignResponse{}, errs.ErrDatabase

			}

			return models.PRReassignResponse{PullRequest: req, ReplacedBy: reviewer.UserID}, nil

		}

	}

	return models.PRReassignResponse{}, errs.ErrNoCandidate

}

// GetPR retrieves all pull requests assigned to a user
func GetPR(ctx context.Context, UserID string) models.UserRequests {

	res, err := database.GetPRFromDBByUser(ctx, UserID)

	if err != nil {

		return models.UserRequests{UserID: UserID, PullRequests: []models.PullRequestShort{}}

	}

	return res

}
