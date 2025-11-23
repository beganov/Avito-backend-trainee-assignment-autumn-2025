package models

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

func UserToTM(user User) TeamMember {
	return TeamMember{
		UserID:   user.UserID,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}
