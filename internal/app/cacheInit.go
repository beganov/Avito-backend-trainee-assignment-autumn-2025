package app

import (
	"context"
	"log"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/cache"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/database"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
)

// LoadCacheFromDB preloads data from database into in-memory cache
func LoadCacheFromDB(ctx context.Context) error {

	cache.InitCache() // Initialize cache instances

	dbCtx, cancel := context.WithTimeout(ctx, config.PostgresTimeOut) // Create context with timeout

	defer cancel()

	rows, err := database.DB.Query(dbCtx, `SELECT team_name FROM teams`) // Load all teams from database into cache

	if err != nil {

		logger.Info("failed to load teams for cache: %v", err)

		return err

	}

	defer rows.Close()

	for rows.Next() {

		var teamName string

		if err := rows.Scan(&teamName); err != nil {

			logger.Info("failed to scan team name: %v", err)

			continue

		}

		team, err, _ := database.GetTeamFromDB(dbCtx, teamName)

		if err != nil {

			logger.Info("failed to load team %s: %v", teamName, err)

			continue

		}

		cache.TeamCache.Set(teamName, team)

	}

	userRows, err := database.DB.Query(dbCtx, `SELECT user_id FROM users`) // Load all users from database into cache

	if err != nil {

		log.Printf("failed to load users for cache: %v", err)

		return err

	}

	defer userRows.Close()

	for userRows.Next() {

		var userID string

		if err := userRows.Scan(&userID); err != nil {

			log.Printf("failed to scan user id: %v", err)

			continue

		}

		user, err, _ := database.GetUserFromDB(dbCtx, userID)

		if err != nil {

			log.Printf("failed to load user %s: %v", userID, err)

			continue

		}

		cache.UserCache.Set(userID, user)

	}

	prRows, err := database.DB.Query(dbCtx, `SELECT pull_request_id FROM pull_requests`) // Load all pr from database into cache

	if err != nil {

		log.Printf("failed to load PRs for cache: %v", err)

		return err

	}

	defer prRows.Close()

	for prRows.Next() {

		var prID string

		if err := prRows.Scan(&prID); err != nil {

			log.Printf("failed to scan PR id: %v", err)

			continue

		}

		pr, err, _ := database.GetPRFromDB(dbCtx, prID)

		if err != nil {

			log.Printf("failed to load PR %s: %v", prID, err)

			continue

		}

		cache.PRcache.Set(prID, pr)

	}

	return nil

}
