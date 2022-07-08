package usecase

import (
	"context"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/domain"
)

type reservationUsecase struct {
	reservationRepo domain.ReservationsRepository
}

func NewReservationUsecase(rv domain.ReservationsRepository) domain.ReservationsUsecase {
	return &reservationUsecase{
		reservationRepo: rv,
	}
}

func (rtu *reservationUsecase) FetchAll(c context.Context) (res []domain.Reservations, err error) {
	res, err = rtu.reservationRepo.FetchAll(c)
	if err != nil {
		return nil, err
	}
	return
}

func (rtu *reservationUsecase) FetchByID(c context.Context, id int) (res domain.Reservations, err error) {
	res, err = rtu.reservationRepo.FetchByID(c, id)
	return
}

func (rtu *reservationUsecase) Store(c context.Context, rt *domain.Reservations) (err error) {
	err = rtu.reservationRepo.Store(c, rt)
	return
}

func (rtu *reservationUsecase) Update(c context.Context, rt *domain.Reservations, id int) (err error) {
	err = rtu.reservationRepo.Update(c, rt, id)
	return
}
