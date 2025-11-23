package database

import (
	"context"
	"fmt"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/models"
)

func GetTeamFromDB(ctx context.Context, teamName string) (models.Team, error, bool) {

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)

	defer cancel()

	var err error

	if DB == nil {

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return models.Team{}, err, false

	}

	var team models.Team

	team.TeamName = teamName

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

func SetTeamToDB(ctx context.Context, team models.Team) error {

	var err error

	if DB == nil {

		err = fmt.Errorf("database not initialized")

		logger.Error(err, err.Error())

		return err

	}

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut)

	defer cancel()

	tx, err := DB.Begin(dbCtx)

	if err != nil {

		logger.Error(err, err.Error())

		return err

	}

	defer tx.Rollback(dbCtx)

	var teamID int

	err = tx.QueryRow(dbCtx, `

        INSERT INTO teams (team_name) 

        VALUES ($1) 

        ON CONFLICT (team_name) DO UPDATE SET team_name = EXCLUDED.team_name

        RETURNING team_id`, team.TeamName).Scan(&teamID)

	if err != nil {

		logger.Error(err, err.Error())

		return err

	}

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

	return tx.Commit(dbCtx)

}
