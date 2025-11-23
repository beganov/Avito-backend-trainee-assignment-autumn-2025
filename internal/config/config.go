package config

import (
	"os"
	"strconv"
	"time"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
)

var (
	PostgresURL string

	TestPostgresURL string

	CacheCap int

	PostgresTimeOut time.Duration

	MigrationPath string
)

func VarsInit() {

	PostgresURL = os.Getenv("POSTGRES_URL")

	TestPostgresURL = os.Getenv("TEST_POSTGRES_URL")

	var err error

	CacheCap, err = strconv.Atoi(os.Getenv("CACHE_CAP"))

	if err != nil {

		logger.Fatal(err, "CACHE_CAP is not number")

	}

	PostgresTimeOutSec, err := strconv.Atoi(os.Getenv("POSTGRES_TIMEOUT"))

	if err != nil {

		logger.Fatal(err, "SELECT_TIMEOUT is not number")

	}

	PostgresTimeOut = time.Duration(PostgresTimeOutSec) * time.Second

	MigrationPath = os.Getenv("MIGRATION_PATH")

}
