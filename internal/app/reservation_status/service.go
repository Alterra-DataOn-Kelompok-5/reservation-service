package locations

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
	ReservationStatusRepository repository.ReservationStatus
}

type Service interface {
	Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.ReservationStatusResponse], error)
	FindByID(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.ReservationStatusResponse, error)
	Store(ctx context.Context, payload *dto.CreateReservationStatusRequestBody) (*dto.ReservationStatusResponse, error)
	UpdateById(ctx context.Context, payload *dto.UpdateReservationStatusRequestBody) (*dto.ReservationStatusResponse, error)
	// DeleteById(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.ReservationStatusWithCUDResponse, error)
}

func NewService(f *factory.Factory) Service {
	return &service{
		ReservationStatusRepository: f.ReservationStatusRepository,
	}
}

func (s *service) Find(ctx context.Context, payload *pkgdto.SearchGetRequest) (*pkgdto.SearchGetResponse[dto.ReservationStatusResponse], error) {
	status, info, err := s.ReservationStatusRepository.FindAll(ctx, payload, &payload.Pagination)
	if err != nil {
		return nil, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	var data []dto.ReservationStatusResponse

	for _, sData := range status {
		data = append(data, dto.ReservationStatusResponse{
			ID:                    sData.ID,
			ReservationStatusName: sData.ReservationStatusName,
		})

	}

	result := new(pkgdto.SearchGetResponse[dto.ReservationStatusResponse])
	result.Data = data
	result.PaginationInfo = *info

	return result, nil
}
func (s *service) FindByID(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.ReservationStatusResponse, error) {
	var result dto.ReservationStatusResponse
	data, err := s.ReservationStatusRepository.FindByID(ctx, payload.ID)
	if err != nil {
		if err == constant.RECORD_NOT_FOUND {
			return &dto.ReservationStatusResponse{}, res.ErrorBuilder(&res.ErrorConstant.NotFound, err)
		}
		return &dto.ReservationStatusResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	result.ID = data.ID
	result.ReservationStatusName = data.ReservationStatusName

	return &result, nil
}

func (s *service) Store(ctx context.Context, payload *dto.CreateReservationStatusRequestBody) (*dto.ReservationStatusResponse, error) {
	var result dto.ReservationStatusResponse
	isExist, err := s.ReservationStatusRepository.ExistByName(ctx, *payload.ReservationStatusName)
	if err != nil {
		return &result, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}
	if isExist {
		return &result, res.ErrorBuilder(&res.ErrorConstant.Duplicate, errors.New("status already exists"))
	}

	data, err := s.ReservationStatusRepository.Save(ctx, payload)
	if err != nil {
		return &result, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	result.ID = data.ID
	result.ReservationStatusName = data.ReservationStatusName

	return &result, nil
}

func (s *service) UpdateById(ctx context.Context, payload *dto.UpdateReservationStatusRequestBody) (*dto.ReservationStatusResponse, error) {
	status, err := s.ReservationStatusRepository.FindByID(ctx, *payload.ID)
	if err != nil {
		if err == constant.RECORD_NOT_FOUND {
			return &dto.ReservationStatusResponse{}, res.ErrorBuilder(&res.ErrorConstant.NotFound, err)
		}
		return &dto.ReservationStatusResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}

	_, err = s.ReservationStatusRepository.Edit(ctx, &status, payload)
	if err != nil {
		return &dto.ReservationStatusResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
	}
	var result dto.ReservationStatusResponse
	result.ID = status.ID
	result.ReservationStatusName = status.ReservationStatusName

	return &result, nil
}

// func (s *service) DeleteById(ctx context.Context, payload *pkgdto.ByIDRequest) (*dto.ReservationStatusWithCUDResponse, error) {
// 	status, err := s.ReservationStatusRepository.FindByID(ctx, payload.ID)
// 	if err != nil {
// 		if err == constant.RECORD_NOT_FOUND {
// 			return &dto.ReservationStatusWithCUDResponse{}, res.ErrorBuilder(&res.ErrorConstant.NotFound, err)
// 		}
// 		return &dto.ReservationStatusWithCUDResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
// 	}
// 	_, err = s.ReservationStatusRepository.Destroy(ctx, &status)
// 	if err != nil {
// 		return &dto.ReservationStatusWithCUDResponse{}, res.ErrorBuilder(&res.ErrorConstant.InternalServerError, err)
// 	}

// 	result := &dto.ReservationStatusWithCUDResponse{
// 		ReservationStatusResponse: dto.ReservationStatusResponse{
// 			ID:                    status.ID,
// 			ReservationStatusName: status.ReservationStatusName,
// 		},
// 		CreatedAt: status.CreatedAt,
// 		UpdatedAt: status.UpdatedAt,
// 		DeletedAt: status.DeletedAt,
// 	}

// 	return result, nil
// }
