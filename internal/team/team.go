package team

import (
	"context"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/database"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

func Get(TeamName string, ctx context.Context) (models.Team, error) {
	resTeam, ok := cache.TeamCache.Get(TeamName)
	if !ok {
		res, err, ok := database.GetTeamFromDB(ctx, TeamName)
		if err != nil {
			return models.Team{}, errs.ErrDatabase
		}
		if !ok {
			return models.Team{}, errs.ErrNotFound
		} else {
			cache.TeamCache.Set(TeamName, res)
			return res, nil
		}
	}
	return resTeam.(models.Team), nil
}

func Add(bindedTeam models.Team, ctx context.Context) (models.TeamResponse, error) {
	_, ok := cache.TeamCache.Get(bindedTeam.TeamName)
	if ok {
		return models.TeamResponse{}, errs.ErrTeamExists
	}
	_, err, ok := database.GetTeamFromDB(ctx, bindedTeam.TeamName)
	if err != nil {
		return models.TeamResponse{}, errs.ErrDatabase
	}
	if ok {
		return models.TeamResponse{}, errs.ErrTeamExists
	}
	cache.TeamCache.Set(bindedTeam.TeamName, bindedTeam)
	err = database.SetTeamToDB(ctx, bindedTeam)
	if err != nil {
		return models.TeamResponse{}, errs.ErrDatabase
	}
	for _, j := range bindedTeam.Members {
		metrics.UsersCreatedTotal.Inc()
		newUser := models.User{
			UserID:   j.UserID,
			Username: j.Username,
			TeamName: bindedTeam.TeamName,
			IsActive: j.IsActive,
		}
		cache.UserCache.Set(j.UserID, newUser)
	}
	return models.TeamResponse{Team: bindedTeam}, nil
}

func SetActive(bindUser models.UserActivity, ctx context.Context) (models.UserResponse, error) {
	var err error
	iUser, ok := cache.UserCache.Get(bindUser.UserID)
	if !ok {
		iUser, err, ok = database.GetUserFromDB(ctx, bindUser.UserID)
		if err != nil {
			return models.UserResponse{}, errs.ErrDatabase
		}
		if !ok {
			return models.UserResponse{}, errs.ErrNoCandidate
		}
	}

	user := iUser.(models.User)
	user.IsActive = bindUser.IsActive
	cache.UserCache.Set(bindUser.UserID, user)

	iTeam, ok := cache.TeamCache.Get(user.TeamName)
	if !ok {
		iTeam, err, ok = database.GetTeamFromDB(ctx, user.TeamName)
		if err != nil || !ok {
			return models.UserResponse{}, errs.ErrDatabase
		}
	}
	team := iTeam.(models.Team)
	for i, j := range team.Members {
		if j.UserID == bindUser.UserID {
			team.Members[i].IsActive = user.IsActive
			break
		}
	}
	cache.TeamCache.Set(user.TeamName, team)
	err = database.SetTeamToDB(ctx, team)
	if err != nil {
		return models.UserResponse{}, errs.ErrDatabase
	}
	return models.UserResponse{User: user}, nil
}
