package domain

import (
	"context"
)

type Reservations struct {
	Model
	EmployeeID         int                  `json:"employee_id"`
	AdminID            int                  `json:"admin_id"`
	ReservationCode    string               `json:"reservation_code"`
	ReservationStatus  int                  `json:"reservation_status" gorm:"default:0"`
	ReservationDesc    string               `json:"reservation_desc" gorm:"type:text"`
	ReservationDetails []ReservationDetails `json:"reservation_detail" gorm:"foreignKey:ReservationID"`
}

type ReservationsUsecase interface {
	FetchAll(ctx context.Context) ([]Reservations, error)
	FetchByID(ctx context.Context, id int) (Reservations, error)
	Store(ctx context.Context, rv *Reservations) error
	Update(ctx context.Context, rv *Reservations, id int) error
	// Delete(ctx context.Context, id int) error
}

type ReservationsRepository interface {
	FetchAll(ctx context.Context) (res []Reservations, err error)
	FetchByID(ctx context.Context, id int) (Reservations, error)
	Store(ctx context.Context, rv *Reservations) error
	Update(ctx context.Context, rv *Reservations, id int) error
	// Delete(ctx context.Context, id int) error
}
