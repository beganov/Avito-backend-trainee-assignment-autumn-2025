package team

import (
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/users"
)

var TeamCache map[string]Team = make(map[string]Team)

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
	if _, ok := TeamCache[TeamName]; !ok {
		return Team{}, errs.ErrNotFound
	}
	return TeamCache[TeamName], nil
}

func Add(bindedTeam Team) (TeamResponse, error) {
	if _, ok := TeamCache[bindedTeam.TeamName]; ok {
		return TeamResponse{}, errs.ErrTeamExists
	}
	TeamCache[bindedTeam.TeamName] = bindedTeam
	for _, j := range bindedTeam.Members {
		users.UserCache[j.UserID] = users.User{
			UserID:   j.UserID,
			Username: j.Username,
			TeamName: bindedTeam.TeamName,
			IsActive: j.IsActive,
		}
	}
	return TeamResponse{bindedTeam}, nil
}
