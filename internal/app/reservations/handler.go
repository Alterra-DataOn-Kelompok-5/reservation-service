package reservations

import (
	"fmt"
	"net/http"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/dto"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/factory"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/pkg/enum"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/pkg/util"
	pkgdto "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/dto"
	res "github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/util/response"
	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) Get(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	_, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	// log.Println(jwtClaims)

	payload := new(pkgdto.SearchGetRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}

	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.Find(c.Request().Context(), payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.CustomSuccessBuilder(http.StatusOK, result.Data, "Get reservations success", &result.PaginationInfo).Send(c)
}

func (h *handler) GetById(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	_, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	// log.Println(jwtClaims)

	payload := new(pkgdto.ByIDRequest)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	result, err := h.service.FindByID(c.Request().Context(), payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

func (h *handler) Create(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	jwtClaims, err := util.ParseJWTToken(authHeader)
	if err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	// log.Println(jwtClaims)

	payload := new(dto.CreateReservationRequestBody)
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}
	payload.EmployeeID = &jwtClaims.UserID

	role, err := h.service.Store(c.Request().Context(), payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(role).Send(c)
}

func (h *handler) UpdateById(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	jwtClaims, err := util.ParseJWTToken(authHeader)
	fmt.Println(jwtClaims.UserID)
	payload := new(dto.UpdateReservationRequestBody)
	payload = &dto.UpdateReservationRequestBody{
		AdminID:             &jwtClaims.UserID,
		ReservationStatusID: payload.ReservationStatusID,
	}
	if err := c.Bind(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
	}

	isAdminOrSameUser := (jwtClaims.UserID == *payload.ID) || (jwtClaims.RoleID == uint(enum.Admin))
	// log.Println(isAdminOrSameUser)
	if (err != nil) || !isAdminOrSameUser {
		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
	}

	result, err := h.service.UpdateById(c.Request().Context(), payload)
	if err != nil {
		return res.ErrorResponse(err).Send(c)
	}

	return res.SuccessResponse(result).Send(c)
}

// func (h *handler) DeleteById(c echo.Context) error {
// 	authHeader := c.Request().Header.Get("Authorization")
// 	jwtClaims, err := util.ParseJWTToken(authHeader)
// 	if (err != nil) || (jwtClaims.RoleID != uint(enum.Admin)) {
// 		return res.ErrorBuilder(&res.ErrorConstant.Unauthorized, err).Send(c)
// 	}

// 	// log.Println(jwtClaims)

// 	payload := new(pkgdto.ByIDRequest)
// 	if err := c.Bind(payload); err != nil {
// 		return res.ErrorBuilder(&res.ErrorConstant.BadRequest, err).Send(c)
// 	}
// 	if err := c.Validate(payload); err != nil {
// 		return res.ErrorBuilder(&res.ErrorConstant.Validation, err).Send(c)
// 	}
// 	result, err := h.service.DeleteById(c.Request().Context(), payload)
// 	if err != nil {
// 		return res.ErrorResponse(err).Send(c)
// 	}

// 	return res.SuccessResponse(result).Send(c)
// }
