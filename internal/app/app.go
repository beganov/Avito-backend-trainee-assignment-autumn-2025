package app

import (
	"github.com/joho/godotenv"

	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/logger"
	"github.com/beganov/Avito-backend-trainee-assignment-autumn-2025/internal/metrics"
)

func Init() {

	if err := godotenv.Load(); err != nil {

		logger.Fatal(err, "No .env file found")

	}

	metrics.Init()

	config.VarsInit()

}
