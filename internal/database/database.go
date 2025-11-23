package database

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
)

var (
	DB *pgxpool.Pool
)

func InitDB(ctx context.Context, dsn string) {

	var err error

	DB, err = pgxpool.New(ctx, dsn)

	if err != nil {

		logger.Fatal(err, "unable to create DB pool")

	}

	if err := DB.Ping(ctx); err != nil {

		logger.Fatal(err, "unable to connect to DB")

	}

}

func RunMigrations(dsn string) {

	DB, err := sql.Open("postgres", dsn)

	if err != nil {

		logger.Fatal(err, "failed to open db for migrations")

	}

	defer DB.Close()

	if err := goose.Up(DB, config.MigrationPath); err != nil {

		logger.Fatal(err, "failed to run migrations")

	}

}
