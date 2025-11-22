package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func InitDB(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err, "unable to create DB pool")
		//logger.Fatal(err, "unable to create DB pool")
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err, "unable to connect to DB")
		//	logger.Fatal(err, "unable to connect to DB")
	}
	return pool
}

// run goose migrations
func RunMigrations(dsn string) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err, "failed to open db for migrations")
		//logger.Fatal(err, "failed to open db for migrations")
	}
	defer db.Close()

	if err := goose.Up(db, config.MigrationPath); err != nil {
		log.Fatal(err, "failed to run migrations")
		//logger.Fatal(err, "failed to run migrations")
	}
}
