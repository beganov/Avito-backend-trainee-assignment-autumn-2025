package models

// User represents a full user entity
type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// UserActivity represents a request to update user activation status
// Used in the setActivity operation
type UserActivity struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// UserResponse is a wrapper for user-related API responses
type UserResponse struct {
	User User `json:"user"`
}

// UserRequests represents a collection of pull requests assigned to a user for review
// Used in the getReview operation
type UserRequests struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
