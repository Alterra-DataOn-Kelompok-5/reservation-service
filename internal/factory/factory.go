package factory

import (
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/repository"
)

type Factory struct {
	ReservationsRepository      repository.Reservations
	ReservationStatusRepository repository.ReservationStatus
}

func NewFactory() *Factory {
	db := database.GetConnection()
	return &Factory{
		repository.NewReservationsRepository(db),
		repository.NewStatusRepository(db),
	}
}
