package team

import (
	"errors"

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

func Get(TeamName string) (Team, error) {
	if _, ok := TeamCache[TeamName]; !ok {
		return Team{}, errors.New("resource not found")
	}
	return TeamCache[TeamName], nil
}

func Add(bindedTeam Team) (Team, error) {
	if _, ok := TeamCache[bindedTeam.TeamName]; ok {
		return Team{}, errors.New("team already exists")
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
	return bindedTeam, nil
}
