package migration

import (
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/domain"
)

var tables = []interface{}{
	&domain.Reservations{},
	&domain.ReservationDetails{},
}

func Migrate() {
	conn := database.GetConnection()
	conn.AutoMigrate(tables...)
}
