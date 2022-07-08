package main

import (
	"log"
	"os"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database/migration"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	errGoEnv := godotenv.Load()
	if errGoEnv != nil {
		// log.Fatal("Error loading .env file")
		panic(errGoEnv)
	}
}

func main() {
	database.CreateConnection()
	migration.Migrate()

	e := echo.New()

	middleware.Init(e)

	log.Fatal(e.Start(":" + os.Getenv("SERVICE_PORT")))
}
