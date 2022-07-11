package model

import "time"

type ReservationStatus struct {
	ID                    uint      `json:"id"`
	CreatedAt             time.Time `json:"created_at" gorm:"type:datetime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"type:datetime"`
	ReservationStatusName string    `json:"reservation_status_name"`
}
