package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

func GetPRFromDBByUser(ctx context.Context, userID string) (models.UserRequests, error) {

	var err error

	if DB == nil {

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return models.UserRequests{}, err

	}

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)

	defer cancel()

	var userRequests models.UserRequests

	userRequests.UserID = userID

	rows, err := DB.Query(dbCtx, `

        SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status

        FROM pull_requests pr

        WHERE pr.assigned_reviewers @> $1`,

		fmt.Sprintf(`["%s"]`, userID))

	if err != nil {

		logger.Error(err, err.Error())

		return models.UserRequests{}, err

	}

	defer rows.Close()

	for rows.Next() {

		var pr models.PullRequestShort

		err := rows.Scan(

			&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status)

		if err != nil {

			logger.Error(err, err.Error())

			return models.UserRequests{}, err

		}

		userRequests.PullRequests = append(userRequests.PullRequests, pr)

	}

	return userRequests, nil

}

func GetUserFromDB(ctx context.Context, userID string) (models.User, error, bool) {

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)

	defer cancel()

	var err error

	if DB == nil {

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return models.User{}, err, false

	}

	var user models.User

	err = DB.QueryRow(dbCtx, `

        SELECT u.user_id, u.username, u.is_active, t.team_name 

        FROM users u 

        JOIN teams t ON u.team_id = t.team_id 

        WHERE u.user_id = $1`, userID).Scan(

		&user.UserID, &user.Username, &user.IsActive, &user.TeamName)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {

			return models.User{}, nil, false

		}

		logger.Error(err, err.Error())

		return models.User{}, err, false

	}

	return user, nil, len(user.UserID) != 0

}
