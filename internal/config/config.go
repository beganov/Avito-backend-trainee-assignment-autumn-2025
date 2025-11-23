package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	PostgresURL string

	CacheCap int

	HttpTimeOut     time.Duration
	PostgresTimeOut time.Duration

	MigrationPath string
)

func VarsInit() {

	PostgresURL = os.Getenv("POSTGRES_URL")

	var err error
	CacheCap, err = strconv.Atoi(os.Getenv("CACHE_CAP"))
	if err != nil {
		log.Fatal(err, "CACHE_CAP is not number")
		//logger.Fatal(err, "CACHE_CAP is not number")
	}

	httpTimeoutSec, err := strconv.Atoi(os.Getenv("HTTP_TIMEOUT"))
	if err != nil {
		log.Fatal(err, "HTTP_TIMEOUT is not number")
		//logger.Fatal(err, "HTTP_TIMEOUT is not number")
	}
	PostgresTimeOutSec, err := strconv.Atoi(os.Getenv("POSTGRES_TIMEOUT"))
	if err != nil {
		log.Fatal(err, "POSTGRES_TIMEOUT is not number")
		//logger.Fatal(err, "SELECT_TIMEOUT is not number")
	}

	HttpTimeOut = time.Duration(httpTimeoutSec) * time.Second
	PostgresTimeOut = time.Duration(PostgresTimeOutSec) * time.Second
	MigrationPath = os.Getenv("MIGRATION_PATH")
}
