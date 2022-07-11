package model

import "time"

type Reservations struct {
	ID                   uint              `json:"id"`
	CreatedAt            time.Time         `json:"created_at" gorm:"type:datetime"`
	UpdatedAt            time.Time         `json:"updated_at" gorm:"type:datetime"`
	ReservationCode      string            `json:"reservation_code"`
	ReservationDesc      string            `json:"reservation_desc" gorm:"type:text"`
	EmployeeID           uint              `json:"employee_id"`
	AdminID              uint              `json:"admin_id" gorm:"default:null"`
	RoomID               uint              `json:"room_id"`
	ReservationTimeStart string            `json:"reservation_time_start" gorm:"type:datetime"`
	ReservationTimeEnd   string            `json:"reservation_time_end" gorm:"type:datetime"`
	ReservationStatusID  uint              `json:"reservation_status_id" gorm:"default:1"`
	ReservationStatus    ReservationStatus `json:"reservation_status" gorm:"foreignKey:ReservationStatusID;references:ID"`
}
