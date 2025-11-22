package users

import (
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/errs"
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

type UserResponse struct {
	user User `json:"user"`
}

func SetActive(bindUser UserActivity) (UserResponse, error) {
	user, ok := UserCache[bindUser.UserID]
	if !ok {
		return UserResponse{}, errs.ErrNotFound
	}
	user.IsActive = bindUser.IsActive
	UserCache[bindUser.UserID] = user
	return UserResponse{user: user}, nil
}
