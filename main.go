package main

import (
	"log"
	"os"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database/migration"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/middleware"

	_reservationHttp "github.com/Alterra-DataOn-Kelompok-5/reservation-service/reservation/delivery/http"
	_reservationRepo "github.com/Alterra-DataOn-Kelompok-5/reservation-service/reservation/repository"
	_reservationUc "github.com/Alterra-DataOn-Kelompok-5/reservation-service/reservation/usecase"

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

	reservationRepo := _reservationRepo.NewMysqlReservationsRepository(database.GetConnection())
	rvu := _reservationUc.NewreservationUsecase(reservationRepo)
	_reservationHttp.NewReservationsHandler(e, rvu)

	log.Fatal(e.Start(":" + os.Getenv("SERVICE_PORT")))
}
