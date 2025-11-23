package database

import (
	"context"
	"fmt"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

// GetTeamFromDB retrieves a team and its members from the database by team name
func GetTeamFromDB(ctx context.Context, teamName string) (models.Team, error, bool) {

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut) // Create context with timeout

	defer cancel()

	var err error

	if DB == nil { // Check if database connection is initialized

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return models.Team{}, err, false

	}

	var team models.Team

	team.TeamName = teamName

	// Query all users belonging to the specified team
	rows, err := DB.Query(dbCtx, `

        SELECT u.user_id, u.username, u.is_active 

        FROM users u 

        JOIN teams t ON u.team_id = t.team_id 

        WHERE t.team_name = $1`, teamName)

	if err != nil {

		logger.Error(err, err.Error())

		return models.Team{}, fmt.Errorf("failed to get team: %w", err), false

	}

	defer rows.Close()

	// Process each team member and add to the team structure
	for rows.Next() {

		var user models.User

		err := rows.Scan(&user.UserID, &user.Username, &user.IsActive)

		if err != nil {

			logger.Error(err, err.Error())

			return models.Team{}, err, false

		}

		user.TeamName = teamName

		team.Members = append(team.Members, models.UserToTM(user))

	}

	return team, nil, len(team.Members) != 0

}

// SetTeamToDB creates or updates a team and all its members in the database
func SetTeamToDB(ctx context.Context, team models.Team) error {

	var err error

	if DB == nil { // Check if database connection is initialized

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return err

	}

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut) // Create context with timeout

	defer cancel()

	tx, err := DB.Begin(dbCtx) // Begin transaction to ensure team creation/update

	if err != nil {

		logger.Error(err, err.Error())

		return err

	}

	defer tx.Rollback(dbCtx) // Ensure rollback if transaction fails

	var teamID int

	// Insert or update team, returning the team_id for user associations
	err = tx.QueryRow(dbCtx, `

        INSERT INTO teams (team_name) 

        VALUES ($1) 

        ON CONFLICT (team_name) DO UPDATE SET team_name = EXCLUDED.team_name

        RETURNING team_id`, team.TeamName).Scan(&teamID)

	if err != nil {

		logger.Error(err, err.Error())

		return err

	}

	// Insert or update each team member
	for _, member := range team.Members {

		_, err := tx.Exec(dbCtx, `

            INSERT INTO users (user_id, username, team_id, is_active) 

            VALUES ($1, $2, $3, $4)

            ON CONFLICT (user_id) DO UPDATE SET 

                username = EXCLUDED.username,

                team_id = EXCLUDED.team_id,

                is_active = EXCLUDED.is_active`,

			member.UserID, member.Username, teamID, member.IsActive)

		if err != nil {

			logger.Error(err, err.Error())

			return err

		}

	}

	// Commit transaction
	return tx.Commit(dbCtx)

}
