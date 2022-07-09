package domain

import (
	"context"
)

type ReservationDetails struct {
	ID                   uint   `gorm:"primaryKey;autoIncrement"`
	RoomID               int    `json:"room_id"`
	ReservationTimeStart string `json:"reservation_time_start" gorm:"type:datetime"`
	ReservationTimeEnd   string `json:"reservation_time_end" gorm:"type:datetime"`
	ReservationID        int    `json:"reservation_id"`
	// Reservations         Reservations `json:"reservations" gorm:"foreignKey:ReservationID;references:ID"`
}

type ReservationDetailsUsecase interface {
	FetchAll(ctx context.Context) ([]ReservationDetails, error)
	FetchByID(ctx context.Context, id int) (ReservationDetails, error)
	Store(ctx context.Context, rvd *[]ReservationDetails) error
	Update(ctx context.Context, rvd *ReservationDetails, id int) error
	FetchByIDandTime(ctx context.Context, id int, start string) (ReservationDetails, error)
	// Delete(ctx context.Context, id int) error
}

type ReservationDetailsRepository interface {
	FetchAll(ctx context.Context) (res []ReservationDetails, err error)
	FetchByID(ctx context.Context, id int) (ReservationDetails, error)
	Store(ctx context.Context, rvd *[]ReservationDetails) error
	Update(ctx context.Context, rvd *ReservationDetails, id int) error
	FetchByIDandTime(ctx context.Context, id int, start string) (ReservationDetails, error)
	// Delete(ctx context.Context, id int) error
}
