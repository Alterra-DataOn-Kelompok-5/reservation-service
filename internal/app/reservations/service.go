package reservations

import (
	"context"
	"errors"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/dto"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/factory"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/repository"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/constant"
	pkgdto "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/dto"
	res "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/util/response"
)

type service struct {
	ReservationsRepository repository.Reservations
}

type Service interface {
	Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.ReservationDetailResponse], error)
	FindByID(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.ReservationDetailResponse, error)
	Store(ctx context.Context, payload *dto.CreateReservationRequestBody) (*dto.ReservationDetailResponse, error)
	UpdateById(ctx context.Context, payload *dto.UpdateReservationRequestBody) (*dto.ReservationDetailResponse, error)
	// DeleteById(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.RoomsWithCUDResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		ReservationsRepository: f.ReservationsRepository,
	}
}

func (s *service) Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.ReservationDetailResponse], error) {
	reservations, info, err := s.ReservationsRepository.FindAll(ctx, payload, &payload.Pagination)
	if err != nil {
		return nil, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	var data []dto.ReservationDetailResponse

	for _, reservation := range reservations {
		data = append(data, dto.ReservationDetailResponse{
			ReservationResponse: dto.ReservationResponse{
				ID:                   reservation.ID,
				ReservationCode:      reservation.ReservationCode,
				ReservationDesc:      reservation.ReservationDesc,
				EmployeeID:           reservation.EmployeeID,
				AdminID:              reservation.AdminID,
				RoomID:               reservation.RoomID,
				ReservationTimeStart: reservation.ReservationTimeStart,
				ReservationTimeEnd:   reservation.ReservationTimeEnd,
			},
			Status: dto.ReservationStatusResponse{
				ID:                    reservation.ReservationStatus.ID,
				ReservationStatusName: reservation.ReservationStatus.ReservationStatusName,
			},
		})

	}

	result := new(pkgdto.SearchGetResponse[dto.ReservationDetailResponse])
	result.Data = data
	result.PaginationInfo = *info

	return result, nil
}
func (s *service) FindByID(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.ReservationDetailResponse, error) {
	data, err := s.ReservationsRepository.FindByID(ctx, payload.ID, true)
	if err != nil {
		if err == constant.RECORD_NOT_FOUND {
			return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.NotFound, err)
		}
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	result := &dto.ReservationDetailResponse{
		ReservationResponse: dto.ReservationResponse{
			ID:                   data.ID,
			ReservationCode:      data.ReservationCode,
			ReservationDesc:      data.ReservationDesc,
			EmployeeID:           data.EmployeeID,
			AdminID:              data.AdminID,
			RoomID:               data.RoomID,
			ReservationTimeStart: data.ReservationTimeStart,
			ReservationTimeEnd:   data.ReservationTimeEnd,
		},
		Status: dto.ReservationStatusResponse{
			ID:                    data.ReservationStatus.ID,
			ReservationStatusName: data.ReservationStatus.ReservationStatusName,
		},
	}

	return result, nil
}

func (s *service) Store(ctx context.Context, payload *dto.CreateReservationRequestBody) (*dto.ReservationDetailResponse, error) {
	isExist, err := s.ReservationsRepository.ExistByCode(ctx, payload.ReservationCode)
	if err != nil {
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}
	if isExist {
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.Duplicate, errors.New("reservation already exists"))
	}

	isRoomReserved, err := s.ReservationsRepository.ExistByReservation(ctx, *payload.RoomID, *payload.ReservationTimeStart)
	if err != nil {
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}
	if isRoomReserved {
		return &dto.ReservationDetailResponse{}, res.CustomErrorBuilder(res.ErrorConstant.Duplicate.Code, res.E_DUPLICATE, "Room already reserved / booked")
	}

	data, err := s.ReservationsRepository.Save(ctx, payload)
	if err != nil {
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	reservation, err := s.ReservationsRepository.FindByID(ctx, data.ID, true)
	if err != nil {
		if err == constant.RECORD_NOT_FOUND {
			return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.NotFound, err)
		}
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	result := &dto.ReservationDetailResponse{
		ReservationResponse: dto.ReservationResponse{
			ID:                   reservation.ID,
			ReservationCode:      reservation.ReservationCode,
			ReservationDesc:      data.ReservationDesc,
			EmployeeID:           reservation.EmployeeID,
			AdminID:              reservation.AdminID,
			RoomID:               reservation.RoomID,
			ReservationTimeStart: reservation.ReservationTimeStart,
			ReservationTimeEnd:   reservation.ReservationTimeEnd,
		},
		Status: dto.ReservationStatusResponse{
			ID:                    reservation.ReservationStatus.ID,
			ReservationStatusName: reservation.ReservationStatus.ReservationStatusName,
		},
	}

	return result, nil
}

func (s *service) UpdateById(ctx context.Context, payload *dto.UpdateReservationRequestBody) (*dto.ReservationDetailResponse, error) {
	reservation, err := s.ReservationsRepository.FindByID(ctx, *payload.ID, false)
	if err != nil {
		if err == constant.RECORD_NOT_FOUND {
			return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.NotFound, err)
		}
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	_, err = s.ReservationsRepository.Edit(ctx, &reservation, payload)
	if err != nil {
		return &dto.ReservationDetailResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	result := &dto.ReservationDetailResponse{
		ReservationResponse: dto.ReservationResponse{
			ID:                   reservation.ID,
			ReservationCode:      reservation.ReservationCode,
			ReservationDesc:      reservation.ReservationDesc,
			EmployeeID:           reservation.EmployeeID,
			AdminID:              reservation.AdminID,
			RoomID:               reservation.RoomID,
			ReservationTimeStart: reservation.ReservationTimeStart,
			ReservationTimeEnd:   reservation.ReservationTimeEnd,
			ReservationStatusID:  reservation.ReservationStatusID,
		},
		Status: dto.ReservationStatusResponse{
			ID:                    reservation.ReservationStatus.ID,
			ReservationStatusName: reservation.ReservationStatus.ReservationStatusName,
		},
	}

	return result, nil
}

// func (s *service) DeleteById(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.RoomsWithCUDResponse, error) {
// 	room, err := s.ReservationsRepository.FindByID(ctx, payload.ID, true)
// 	if err != nil {
// 		if err == constant.RECORD_NOT_FOUND {
// 			return &dto.RoomsWithCUDResponse{}, res.ErrorBuilder(&res.ErrorConstant.NotFound, err)
// 		}
// 		return &dto.RoomsWithCUDResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
// 	}
// 	_, err = s.ReservationsRepository.Destroy(ctx, &room)
// 	if err != nil {
// 		return &dto.RoomsWithCUDResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
// 	}

// 	result := &dto.RoomsWithCUDResponse{
// 		RoomsResponse: dto.RoomsResponse{
// 			ID:       room.ID,
// 			RoomName: room.RoomName,
// 			RoomDesc: room.RoomDesc,
// 		},
// 		CreatedAt: room.CreatedAt,
// 		UpdatedAt: room.UpdatedAt,
// 		DeletedAt: room.DeletedAt,
// 	}

// 	return result, nil
// }
