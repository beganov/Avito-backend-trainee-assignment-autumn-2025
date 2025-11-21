package users

import (
	"errors"

	pullrequest "github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/pullRequest"
)

var UserCache map[string]User = make(map[string]User)

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserActivity struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserRequests struct {
	UserID        string                         `json:"user_id"`
	Pull_requests []pullrequest.PullRequestShort `json:"pull_requests"`
}

func GetPR(UserID string) UserRequests {

	return UserRequests{}
}

func SetActive(bindUser UserActivity) (User, error) {
	if _, ok := UserCache[bindUser.UserID]; !ok {
		return User{}, errors.New("resource not found")
	}
	user := UserCache[bindUser.UserID]
	user.IsActive = bindUser.IsActive
	UserCache[bindUser.UserID] = user
	return user, nil
}
