package models

// Team represents a team
type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// TeamMember represents a user within a team context
type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// TeamResponse is a wrapper for team-related API responses
type TeamResponse struct {
	Team Team `json:"team"`
}

// UserToTM converts a full User model to a TeamMember representation
func UserToTM(user User) TeamMember {
	return TeamMember{
		UserID:   user.UserID,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}
