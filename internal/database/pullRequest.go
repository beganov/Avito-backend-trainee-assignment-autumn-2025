package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

func GetPRFromDB(ctx context.Context, prID string) (models.PullRequest, error, bool) {

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)

	defer cancel()

	var err error

	if DB == nil {

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return models.PullRequest{}, err, false

	}

	var pr models.PullRequest

	var createdAt, mergedAt sql.NullTime

	err = DB.QueryRow(dbCtx, `

        SELECT pull_request_id, pull_request_name, author_id, status, 

               assigned_reviewers, created_at, merged_at

        FROM pull_requests 

        WHERE pull_request_id = $1`, prID).Scan(

		&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status,

		&pr.AssignedReviewers, &createdAt, &mergedAt)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {

			return models.PullRequest{}, nil, false

		}

		logger.Error(err, err.Error())

		return models.PullRequest{}, err, false

	}

	if createdAt.Valid {

		pr.CreatedAt = createdAt.Time.Format(time.RFC3339)

	}

	if mergedAt.Valid {

		pr.MergedAt = mergedAt.Time.Format(time.RFC3339)

	}

	return pr, nil, len(pr.PullRequestID) != 0

}

func SetPRToDB(ctx context.Context, pr models.PullRequest) error {

	var err error

	if DB == nil {

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return err

	}

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)

	defer cancel()

	var createdAt, mergedAt interface{}

	if pr.CreatedAt != "" {

		t, err := time.Parse(time.RFC3339, pr.CreatedAt)

		if err != nil {

			logger.Error(err, err.Error())

			return err

		}

		createdAt = t

	} else {

		createdAt = nil

	}

	if pr.MergedAt != "" {

		t, err := time.Parse(time.RFC3339, pr.MergedAt)

		if err != nil {

			logger.Error(err, err.Error())

			return err

		}

		mergedAt = t

	} else {

		mergedAt = nil

	}

	_, err = DB.Exec(dbCtx, `

        INSERT INTO pull_requests 

        (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at)

        VALUES ($1, $2, $3, $4, $5, $6, $7)

        ON CONFLICT (pull_request_id) DO UPDATE SET

            pull_request_name = EXCLUDED.pull_request_name,

            author_id = EXCLUDED.author_id,

            status = EXCLUDED.status,

            assigned_reviewers = EXCLUDED.assigned_reviewers,

            merged_at = EXCLUDED.merged_at`,

		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status,

		pr.AssignedReviewers, createdAt, mergedAt)

	if err != nil {

		logger.Error(err, err.Error())

		return err

	}

	return nil

}
