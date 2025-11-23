package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	DB *pgxpool.Pool
)

func InitDB(ctx context.Context, dsn string) {
	var err error
	DB, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err, "unable to create DB pool")
		//logger.Fatal(err, "unable to create DB pool")
	}
	if err := DB.Ping(ctx); err != nil {
		log.Fatal(err, "unable to connect to DB")
		//	logger.Fatal(err, "unable to connect to DB")
	}
}

// run goose migrations
func RunMigrations(dsn string) {
	DB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err, "failed to open db for migrations")
		//logger.Fatal(err, "failed to open db for migrations")
	}
	defer DB.Close()

	if err := goose.Up(DB, config.MigrationPath); err != nil {
		log.Fatal(err, "failed to run migrations")
		//logger.Fatal(err, "failed to run migrations")
	}
}

type iTeam interface {
}

type iUser interface {
}

type iPR interface {
}

func GetTeamFromDB(ctx context.Context, teamName string) (models.Team, error, bool) {
	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)
	defer cancel()

	if DB == nil {
		return models.Team{}, fmt.Errorf("database not initialized"), false
	}

	var team models.Team
	team.TeamName = teamName

	rows, err := DB.Query(dbCtx, `
        SELECT u.user_id, u.username, u.is_active 
        FROM users u 
        JOIN teams t ON u.team_id = t.team_id 
        WHERE t.team_name = $1`, teamName)
	if err != nil {
		return models.Team{}, fmt.Errorf("failed to get team: %w", err), false
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.IsActive)
		if err != nil {
			return models.Team{}, fmt.Errorf("failed to scan user: %w", err), false
		}
		user.TeamName = teamName
		team.Members = append(team.Members, models.UserToTM(user))
	}
	return team, nil, len(team.Members) != 0
}

func GetUserFromDB(ctx context.Context, userID string) (models.User, error, bool) {
	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)
	defer cancel()

	if DB == nil {
		return models.User{}, fmt.Errorf("database not initialized"), false
	}

	var user models.User
	err := DB.QueryRow(dbCtx, `
        SELECT u.user_id, u.username, u.is_active, t.team_name 
        FROM users u 
        JOIN teams t ON u.team_id = t.team_id 
        WHERE u.user_id = $1`, userID).Scan(
		&user.UserID, &user.Username, &user.IsActive, &user.TeamName)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, nil, false
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err), false
	}
	return user, nil, len(user.UserID) != 0
}

func GetPRFromDB(ctx context.Context, prID string) (models.PullRequest, error, bool) {
	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)
	defer cancel()

	if DB == nil {
		return models.PullRequest{}, fmt.Errorf("database not initialized"), false
	}

	var pr models.PullRequest
	var createdAt, mergedAt sql.NullTime

	err := DB.QueryRow(dbCtx, `
        SELECT pull_request_id, pull_request_name, author_id, status, 
               assigned_reviewers, created_at, merged_at
        FROM pull_requests 
        WHERE pull_request_id = $1`, prID).Scan(
		&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status,
		&pr.AssignedReviewers, &createdAt, &mergedAt)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, pgx.ErrNoRows) {
			return models.PullRequest{}, nil, false
		}
		return models.PullRequest{}, fmt.Errorf("failed to get PR: %w", err), false
	}

	if createdAt.Valid {
		pr.CreatedAt = createdAt.Time.Format(time.RFC3339)
	}
	if mergedAt.Valid {
		pr.MergedAt = mergedAt.Time.Format(time.RFC3339)
	}

	return pr, nil, len(pr.PullRequestID) != 0
}

func SetTeamToDB(ctx context.Context, team models.Team) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	tx, err := DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var teamID int
	err = tx.QueryRow(ctx, `
        INSERT INTO teams (team_name) 
        VALUES ($1) 
        ON CONFLICT (team_name) DO UPDATE SET team_name = EXCLUDED.team_name
        RETURNING team_id`, team.TeamName).Scan(&teamID)
	if err != nil {
		return fmt.Errorf("failed to insert team: %w", err)
	}

	for _, member := range team.Members {
		_, err := tx.Exec(ctx, `
            INSERT INTO users (user_id, username, team_id, is_active) 
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id) DO UPDATE SET 
                username = EXCLUDED.username,
                team_id = EXCLUDED.team_id,
                is_active = EXCLUDED.is_active`,
			member.UserID, member.Username, teamID, member.IsActive)
		if err != nil {
			return fmt.Errorf("failed to upsert user %s: %w", member.UserID, err)
		}
	}

	return tx.Commit(ctx)
}

func SetUserToDB(iUser) error {
	return errors.New("error")
}

func SetPRToDB(ctx context.Context, pr models.PullRequest) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	var createdAt, mergedAt interface{}

	if pr.CreatedAt != "" {
		t, err := time.Parse(time.RFC3339, pr.CreatedAt)
		if err != nil {
			return fmt.Errorf("invalid created_at format: %w", err)
		}
		createdAt = t
	} else {
		createdAt = nil // или time.Now() если хочешь текущее время
	}

	if pr.MergedAt != "" {
		t, err := time.Parse(time.RFC3339, pr.MergedAt)
		if err != nil {
			return fmt.Errorf("invalid merged_at format: %w", err)
		}
		mergedAt = t
	} else {
		mergedAt = nil
	}

	_, err := DB.Exec(ctx, `
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
		fmt.Println(err)
		return fmt.Errorf("failed to save PR: %w", err)
	}

	return nil
}
