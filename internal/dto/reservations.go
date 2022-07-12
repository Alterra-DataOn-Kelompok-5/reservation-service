package dto

import (
	"time"
)

type (
	CreateReservationRequestBody struct {
		ReservationCode      string `json:"reservation_code" validate:"required"`
		ReservationDesc      *string `json:"reservation_desc" validate:"required"`
		EmployeeID           *uint   `json:"employee_id" validate:"required"`
		RoomID               *uint   `json:"room_id" validate:"required"`
		ReservationTimeStart *string `json:"reservation_time_start" validate:"required"`
		ReservationTimeEnd   *string `json:"reservation_time_end" validate:"required"`
	}
	UpdateReservationRequestBody struct {
		ID                  *uint `param:"id" validate:"required"`
		AdminID             *uint `json:"admin_id" validate:"omitempty"`
		ReservationStatusID *uint `json:"reservation_status_id" validate:"omitempty"`
	}
	ReservationResponse struct {
		ID                   uint   `json:"id"`
		ReservationCode      string `json:"reservation_code"`
		ReservationDesc      string `json:"reservation_desc"`
		EmployeeID           uint   `json:"employee_id"`
		AdminID              uint   `json:"admin_id"`
		RoomID               uint   `json:"room_id"`
		ReservationTimeStart string `json:"reservation_time_start"`
		ReservationTimeEnd   string `json:"reservation_time_end"`
		ReservationStatusID  uint   `json:"reservation_status_id"`
	}
	ReservationWithCUDResponse struct {
		ReservationResponse
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	ReservationDetailResponse struct {
		ReservationResponse
		Status ReservationStatusResponse `json:"status"`
	}
)
