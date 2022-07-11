package seeder

import (
	"log"
	"time"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/model"
	"gorm.io/gorm"
)

func reservationSeeder(db *gorm.DB) {
	now := time.Now()
	var reservation = []model.Reservations{
		{
			ID:                   1,
			CreatedAt:            now,
			UpdatedAt:            now,
			EmployeeID:           2,
			AdminID:              0,
			ReservationCode:      "RSVN/20220710/001",
			ReservationDesc:      "Reservasi ruang meeting",
			RoomID:               1,
			ReservationTimeStart: "2022-07-10 10:00:00",
			ReservationTimeEnd:   "2022-07-10 11:00:00",
			ReservationStatusID:  1,
		},
		{
			ID:                   2,
			CreatedAt:            now,
			UpdatedAt:            now,
			EmployeeID:           3,
			AdminID:              1,
			ReservationCode:      "RSVN/20220710/002",
			ReservationDesc:      "Reservasi ruang meeting",
			RoomID:               2,
			ReservationTimeStart: "2022-07-10 12:00:00",
			ReservationTimeEnd:   "2022-07-10 13:00:00",
			ReservationStatusID:  2,
		},
		{
			ID:                   3,
			CreatedAt:            now,
			UpdatedAt:            now,
			EmployeeID:           2,
			AdminID:              1,
			ReservationCode:      "RSVN/20220710/003",
			ReservationDesc:      "Reservasi ruang meeting",
			RoomID:               3,
			ReservationTimeStart: "2022-07-10 14:00:00",
			ReservationTimeEnd:   "2022-07-10 15:00:00",
			ReservationStatusID:  3,
		},
	}
	if err := db.Create(&reservation).Error; err != nil {
		log.Printf("cannot seed data reservation, with error %v\n", err)
	}
	log.Println("success seed data reservation")
}
