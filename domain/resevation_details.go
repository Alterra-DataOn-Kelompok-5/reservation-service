package domain

import (
	"context"
	"time"
)

type ReservationDetails struct {
	ID                   uint         `gorm:"primaryKey;autoIncrement"`
	RoomID               int          `json:"room_id"`
	ReservationTimeStart time.Time    `json:"reservation_time_start"`
	ReservationTimeEnd   time.Time    `json:"reservation_time_end"`
	ReservationID        int          `json:"reservation_id"`
	Reservations         Reservations `json:"reservations" gorm:"foreignKey:ReservationID;references:ID"`
}

type ReservationDetailsUsecase interface {
	FetchAll(ctx context.Context) ([]ReservationDetails, error)
	FetchByID(ctx context.Context, id int) (ReservationDetails, error)
	Store(ctx context.Context, rvd *ReservationDetails) error
	Update(ctx context.Context, rvd *ReservationDetails, id int) error
	// Delete(ctx context.Context, id int) error
}

type ReservationDetailsRepository interface {
	FetchAll(ctx context.Context) (res []ReservationDetails, err error)
	FetchByID(ctx context.Context, id int) (ReservationDetails, error)
	Store(ctx context.Context, rvd *ReservationDetails) error
	Update(ctx context.Context, rvd *ReservationDetails, id int) error
	// Delete(ctx context.Context, id int) error
}
