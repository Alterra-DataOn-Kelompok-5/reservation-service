package seeder

import (
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/database"
	"gorm.io/gorm"
)

type seed struct {
	DB *gorm.DB
}

func NewSeeder() *seed {
	return &seed{database.GetConnection()}
}

func (s *seed) SeedAll() {
	reservationStatusSeeder(s.DB)
	reservationSeeder(s.DB)
}

func (s *seed) DeleteAll() {
	s.DB.Exec("DELETE FROM reservations")
	s.DB.Exec("DELETE FROM reservation_statuses")
}
