package dto

import (
	"time"
)

type (
	CreateReservationStatusRequestBody struct {
		ReservationStatusName *string `json:"reservation_status_name" validate:"required"`
	}
	UpdateReservationStatusRequestBody struct {
		ID                    *uint   `param:"id" validate:"required"`
		ReservationStatusName *string `json:"reservation_status_name" validate:"required"`
	}
	ReservationStatusResponse struct {
		ID                    uint   `json:"id"`
		ReservationStatusName string `json:"reservation_status_name"`
	}
	ReservationStatusWithCUDResponse struct {
		ReservationStatusResponse
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
