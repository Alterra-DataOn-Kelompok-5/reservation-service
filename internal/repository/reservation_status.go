package repository

import (
	"context"
	"strings"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/dto"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/model"
	pkgdto "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/dto"
	"gorm.io/gorm"
)

type ReservationStatus interface {
	FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.ReservationStatus, *pkgdto.PaginationInfo, error)
	FindByID(ctx context.Context, id uint) (model.ReservationStatus, error)
	Save(ctx context.Context, status *dto.CreateReservationStatusRequestBody) (model.ReservationStatus, error)
	Edit(ctx context.Context, oldStatus *model.ReservationStatus, updateData *dto.UpdateReservationStatusRequestBody) (*model.ReservationStatus, error)
	// Destroy(ctx context.Context, status *model.ReservationStatus) (*model.ReservationStatus, error)
	ExistByName(ctx context.Context, name string) (bool, error)
	ExistByID(ctx context.Context, id uint) (bool, error)
}

type reservationStatus struct {
	Db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) *reservationStatus {
	return &reservationStatus{
		db,
	}
}

func (r *reservationStatus) FindAll(ctx context.Context, payload *pkgdto.SearchGetRequest, pagination *pkgdto.Pagination) ([]model.ReservationStatus, *pkgdto.PaginationInfo, error) {
	var status []model.ReservationStatus
	var count int64

	query := r.Db.WithContext(ctx).Model(&model.ReservationStatus{})

	if payload.Search != "" {
		search := "%" + strings.ToLower(payload.Search) + "%"
		query = query.Where("lower(reservation_status_name) LIKE ?", search, search)
	}

	countQuery := query
	if err := countQuery.Count(&count).Error; err != nil {
		return nil, nil, err
	}

	limit, offset := pkgdto.GetLimitOffset(pagination)

	err := query.Limit(limit).Offset(offset).Find(&status).Error

	return status, pkgdto.CheckInfoPagination(pagination, count), err
}

func (r *reservationStatus) FindByID(ctx context.Context, id uint) (model.ReservationStatus, error) {
	var status model.ReservationStatus
	if err := r.Db.WithContext(ctx).Model(&model.ReservationStatus{}).Where("id = ?", id).First(&status).Error; err != nil {
		return status, err
	}
	return status, nil
}

func (r *reservationStatus) Save(ctx context.Context, status *dto.CreateReservationStatusRequestBody) (model.ReservationStatus, error) {
	newStatus := model.ReservationStatus{
		ReservationStatusName: *status.ReservationStatusName,
	}
	if err := r.Db.WithContext(ctx).Save(&newStatus).Error; err != nil {
		return newStatus, err
	}
	return newStatus, nil
}

func (r *reservationStatus) Edit(ctx context.Context, oldStatus *model.ReservationStatus, updateData *dto.UpdateReservationStatusRequestBody) (*model.ReservationStatus, error) {
	if updateData.ReservationStatusName != nil {
		oldStatus.ReservationStatusName = *updateData.ReservationStatusName
	}

	if err := r.Db.WithContext(ctx).Save(oldStatus).Find(oldStatus).Error; err != nil {
		return nil, err
	}

	return oldStatus, nil
}

// func (r *reservationStatus) Destroy(ctx context.Context, status *model.ReservationStatus) (*model.ReservationStatus, error) {
// 	if err := r.Db.WithContext(ctx).Delete(status).Error; err != nil {
// 		return nil, err
// 	}
// 	return status, nil
// }

func (r *reservationStatus) ExistByName(ctx context.Context, name string) (bool, error) {
	var (
		count   int64
		isExist bool
	)
	if err := r.Db.WithContext(ctx).Model(&model.ReservationStatus{}).Where("reservation_status_name = ?", name).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}

func (r *reservationStatus) ExistByID(ctx context.Context, id uint) (bool, error) {
	var (
		count   int64
		isExist bool
	)
	if err := r.Db.WithContext(ctx).Model(&model.ReservationStatus{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return isExist, err
	}
	if count > 0 {
		isExist = true
	}
	return isExist, nil
}
