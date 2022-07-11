package seeder

import (
	"log"
	"time"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/model"
	"gorm.io/gorm"
)

func reservationStatusSeeder(db *gorm.DB) {
	now := time.Now()
	var status = []model.ReservationStatus{
		{
			ID:                    1,
			CreatedAt:             now,
			UpdatedAt:             now,
			ReservationStatusName: "Reserved",
		},
		{
			ID:                    2,
			CreatedAt:             now,
			UpdatedAt:             now,
			ReservationStatusName: "Booked",
		},
		{
			ID:                    3,
			CreatedAt:             now,
			UpdatedAt:             now,
			ReservationStatusName: "Completed",
		},
		{
			ID:                    4,
			CreatedAt:             now,
			UpdatedAt:             now,
			ReservationStatusName: "Cancelled",
		},
	}
	if err := db.Create(&status).Error; err != nil {
		log.Printf("cannot seed data status, with error %v\n", err)
	}
	log.Println("success seed data status")
}
