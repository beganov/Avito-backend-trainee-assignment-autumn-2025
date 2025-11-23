package pullrequest

import (
	"context"
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/database"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

//var UserPRcache map[string][]models.PullRequestShort = make(map[string][]models.PullRequestShort)

func Create(ctx context.Context, bindedPR models.PullRequestShort) (models.PRResponse, error) {
	_, ok := cache.PRcache.Get(bindedPR.PullRequestID)
	if ok {
		return models.PRResponse{}, errs.ErrPRExists
	}
	_, err, ok := database.GetPRFromDB(ctx, bindedPR.PullRequestID)
	if err != nil {
		return models.PRResponse{}, errs.ErrNoCandidate //вернуть 500
	}
	if ok {
		return models.PRResponse{}, errs.ErrPRExists
	}
	iUser, ok := cache.UserCache.Get(bindedPR.AuthorID)
	if !ok {
		iUser, err, ok = database.GetUserFromDB(ctx, bindedPR.AuthorID)
		if err != nil {
			return models.PRResponse{}, errs.ErrPRExists //вернуть 500
		}
		if !ok {
			return models.PRResponse{}, errs.ErrNotFound
		}
	}
	author := iUser.(models.User)
	req := models.PullRequest{
		PullRequestID:     bindedPR.PullRequestID,
		PullRequestName:   bindedPR.PullRequestName,
		AuthorID:          bindedPR.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: []string{},
		CreatedAt:         time.Now().UTC().Format(time.RFC3339),
	}
	iTeam, ok := cache.TeamCache.Get(author.TeamName)
	if !ok {
		iTeam, err, ok = database.GetTeamFromDB(ctx, author.TeamName)
		if err != nil || !ok {
			return models.PRResponse{}, errs.ErrPRExists //вернуть 500
		}
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
		return models.PRResponse{}, errs.ErrPRExists //вернуть 500
	}
	return models.PRResponse{PullRequest: req}, nil
}

func Merge(ctx context.Context, bindedPR models.PullRequestShort) (models.PRResponse, error) {
	var err error
	iPR, ok := cache.PRcache.Get(bindedPR.PullRequestID)
	if !ok {
		iPR, err, ok = database.GetPRFromDB(ctx, bindedPR.PullRequestID)
		if err != nil {
			return models.PRResponse{}, errs.ErrNoCandidate //вернуть 500
		}
		if !ok {
			return models.PRResponse{}, errs.ErrNotFound
		}
	}
	req := iPR.(models.PullRequest)
	if req.Status == "MERGED" {
		return models.PRResponse{PullRequest: req}, nil
	}
	req.Status = "MERGED"
	req.MergedAt = time.Now().UTC().Format(time.RFC3339)
	cache.PRcache.Set(bindedPR.PullRequestID, req)
	err = database.SetPRToDB(ctx, req)
	if err != nil {
		return models.PRResponse{}, errs.ErrPRExists //вернуть 500
	}
	return models.PRResponse{PullRequest: req}, nil
}

func Reassign(ctx context.Context, bindedPR models.PRReassign) (models.PRReassignResponse, error) {
	var err error
	iPR, ok := cache.PRcache.Get(bindedPR.PullRequestID)
	if !ok {
		iPR, err, ok = database.GetPRFromDB(ctx, bindedPR.PullRequestID)
		if err != nil {
			return models.PRReassignResponse{}, errs.ErrNoCandidate //вернуть 500
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
			return models.PRReassignResponse{}, errs.ErrNoCandidate //вернуть 500
		}
		if !ok {
			return models.PRReassignResponse{}, errs.ErrNotFound
		}
	}
	reviewer := iUser.(models.User)

	if req.Status == "MERGED" {
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
			return models.PRReassignResponse{}, errs.ErrPRExists //вернуть 500
		}
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
				return models.PRReassignResponse{}, errs.ErrPRExists //вернуть 500
			}
			return models.PRReassignResponse{PullRequest: req, ReplacedBy: reviewer.UserID}, nil
		}
	}
	return models.PRReassignResponse{}, errs.ErrNoCandidate
}

func GetPR(ctx context.Context, UserID string) models.UserRequests {
	res, err := database.GetPRFromDBByUser(ctx, UserID)
	if err != nil {
		return models.UserRequests{UserID: UserID, PullRequests: []models.PullRequestShort{}} //вернуть 500
	}
	return res
}

func PrToPrShort(PR models.PullRequest) models.PullRequestShort {
	return models.PullRequestShort{
		PullRequestID:   PR.PullRequestID,
		PullRequestName: PR.PullRequestName,
		AuthorID:        PR.AuthorID,
		Status:          PR.Status,
	}
}
