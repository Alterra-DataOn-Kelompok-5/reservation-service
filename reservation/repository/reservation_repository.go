package repository

import (
	"context"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/domain"
	"gorm.io/gorm"
)

type mysqlReservationsRepository struct {
	Db *gorm.DB
}

func NewMysqlReservationsRepository(conn *gorm.DB) domain.ReservationsRepository {
	return &mysqlReservationsRepository{conn}
}

func (m *mysqlReservationsRepository) FetchAll(ctx context.Context) (res []domain.Reservations, err error) {
	var reservations []domain.Reservations
	err = m.Db.WithContext(ctx).Model(&domain.Reservations{}).Find(&reservations).Error

	return reservations, err
}

func (m *mysqlReservationsRepository) FetchByID(ctx context.Context, id int) (res domain.Reservations, err error) {
	var reservation domain.Reservations
	err = m.Db.WithContext(ctx).Model(&domain.Reservations{}).Where("id = ?", id).First(&reservation).Error
	return reservation, err
}

func (m *mysqlReservationsRepository) Store(ctx context.Context, rv *domain.Reservations) error {
	err := m.Db.WithContext(ctx).Model(&domain.Reservations{}).Create(&rv).Error

	return err
}

func (m *mysqlReservationsRepository) Update(ctx context.Context, rv *domain.Reservations, id int) error {
	err := m.Db.WithContext(ctx).Model(&domain.Reservations{}).Where("id = ?", id).Updates(&rv).Error

	return err
}
