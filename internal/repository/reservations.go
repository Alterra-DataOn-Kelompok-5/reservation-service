package repository

import (
	"context"
	"strings"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/dto"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/model"
	pkgdto "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/dto"
	"gorm.io/gorm"
)

type Reservations interface {
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, p *pkgdto.Pagination) ([]model.Reservations, *pkgdto.PaginationInfo, error)
	FindByID(ctx context.Context, id uint, usePreload bool) (model.Reservations, error)
	ExistByCode(ctx context.Context, code string) (bool, error)
	ExistByReservation(ctx context.Context, roomID uint, startTime string) (bool, error)
	Save(ctx context.Context, reservation *dto.CreateReservationRequestBody) (model.Reservations, error)
	Edit(ctx context.Context, oldReservation *model.Reservations, updateData *dto.UpdateReservationRequestBody) (*model.Reservations, error)
	// Destroy(ctx context.Context, reservation *model.Reservations) (*model.Reservations, error)
}

type reservations struct {
	Db *gorm.DB
}

func NewReservationsRepository(db *gorm.DB) *reservations {
	return &reservations{
		db,
	}
}

func (r *reservations) FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.Reservations, *pkgdto.PaginationInfo, error) {
	var reservations []model.Reservations
	var count int64

	query := r.Db.WithContext(ctx).Model(&model.Reservations{}).Preload("ReservationStatus")

	if payload.Search != "" {
		search := "%" + strings.ToLower(payload.Search) + "%"
		query = query.Where("lower(reservation_code) LIKE ?", search, search)
	}

	countQuery := query
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, nil, err
	}

	limit, offset := pkgdto.GetLimitOffset(pagination)

	err := query.Limit(limit).Offset(offset).Find(&reservations).Error

	return reservations, pkgdto.CheckInfoPagination(pagination, count), err
}

func (r *reservations) FindByID(ctx context.Context, id uint, usePreload bool) (model.Reservations, error) {
	var reservation model.Reservations
	q := r.Db.WithContext(ctx).Model(&model.Reservations{}).Where("id = ?", id)
	if usePreload {
		q = q.Preload("ReservationStatus")
	}
	err := q.First(&reservation).Error
	return reservation, err
}

func (r *reservations) ExistByCode(ctx context.Context, code string) (bool, error) {
	var (
		count   int64
		isExist bool
	)
	if err := r.Db.WithContext(ctx).Model(&model.Reservations{}).Where("reservation_code = ?", code).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}

func (r *reservations) ExistByReservation(ctx context.Context, roomID uint, startTime string) (bool, error) {
	var (
		count   int64
		isExist bool
	)
	err := r.Db.WithContext(ctx).Model(&model.Reservations{}).
		Where("room_id = ? and ? between reservation_time_start and reservation_time_end and (reservation_status_id <> 1 or reservation_status_id <> 2)", roomID, startTime).
		Count(&count).Error

	if err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}

func (r *reservations) Save(ctx context.Context, reservation *dto.CreateReservationRequestBody) (model.Reservations, error) {
	var newReservation model.Reservations
	newReservation.ReservationCode = reservation.ReservationCode
	if reservation.ReservationDesc != nil {
		newReservation.ReservationDesc = *reservation.ReservationDesc
	}
	if reservation.EmployeeID != nil {
		newReservation.EmployeeID = *reservation.EmployeeID
	}
	if reservation.RoomID != nil {
		newReservation.RoomID = *reservation.RoomID
	}
	if reservation.ReservationTimeStart != nil {
		newReservation.ReservationTimeStart = *reservation.ReservationTimeStart
	}
	if reservation.ReservationTimeEnd != nil {
		newReservation.ReservationTimeEnd = *reservation.ReservationTimeEnd
	}
	if err := r.Db.WithContext(ctx).Save(&newReservation).Error; err != nil {
		return newReservation, err
	}
	return newReservation, nil
}

func (r *reservations) Edit(ctx context.Context, oldReservation *model.Reservations, updateData *dto.UpdateReservationRequestBody) (*model.Reservations, error) {
	if updateData.AdminID != nil {
		oldReservation.AdminID = *updateData.AdminID
	}
	if updateData.ReservationStatusID != nil {
		oldReservation.ReservationStatusID = *updateData.ReservationStatusID
	}

	if err := r.Db.
		WithContext(ctx).
		Save(oldReservation).
		Preload("ReservationStatus").
		Find(oldReservation).
		Error; err != nil {
		return nil, err
	}

	return oldReservation, nil
}

func (r *reservations) Destroy(ctx context.Context, reservation *model.Reservations) (*model.Reservations, error) {
	if err := r.Db.WithContext(ctx).Delete(reservation).Error; err != nil {
		return nil, err
	}
	return reservation, nil
}
