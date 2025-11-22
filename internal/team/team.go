package team

import (
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/users"
)

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamResponse struct {
	Team Team `json:"team"`
}

func Get(TeamName string) (Team, error) {
	resTeam, ok := cache.TeamCache.Get(TeamName)
	if !ok {
		return Team{}, errs.ErrNotFound
	}
	return resTeam.(Team), nil
}

func Add(bindedTeam Team) (TeamResponse, error) {
	_, ok := cache.TeamCache.Get(bindedTeam.TeamName)
	if ok {
		return TeamResponse{}, errs.ErrTeamExists
	}
	cache.TeamCache.Set(bindedTeam.TeamName, bindedTeam)
	for _, j := range bindedTeam.Members {
		cache.UserCache.Set(j.UserID, users.User{
			UserID:   j.UserID,
			Username: j.Username,
			TeamName: bindedTeam.TeamName,
			IsActive: j.IsActive,
		})
	}
	return TeamResponse{bindedTeam}, nil
}

func SetActive(bindUser UserActivity) (UserResponse, error) {
	iUser, ok := cache.UserCache.Get(bindUser.UserID)
	if !ok {
		return UserResponse{}, errs.ErrNotFound
	}
	user := iUser.(users.User)
	user.IsActive = bindUser.IsActive
	cache.UserCache.Set(bindUser.UserID, user)
	iTeam, _ := cache.TeamCache.Get(user.TeamName)
	team := iTeam.(Team)
	for i, j := range team.Members {
		if j.UserID == bindUser.UserID {
			team.Members[i].IsActive = user.IsActive
			break
		}
	}
	cache.TeamCache.Set(user.TeamName, team)
	return UserResponse{User: user}, nil
}

type UserActivity struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User users.User `json:"user"`
}
