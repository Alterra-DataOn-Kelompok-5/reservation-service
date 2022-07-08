package http

import (
	"net/http"
	"strconv"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/domain"
	"github.com/labstack/echo/v4"
)

type ReservationsHandler struct {
	ReservationsUsecase domain.ReservationsUsecase
}

func NewReservationsHandler(e *echo.Echo, rvu domain.ReservationsUsecase) {
	handler := &ReservationsHandler{
		ReservationsUsecase: rvu,
	}

	e.GET("/reservations", handler.FetchAllReservations)
	e.GET("/reservations/:id", handler.FetchReservationByID)
	e.POST("/reservations", handler.CreateReservation)
	e.PUT("/reservations/:id", handler.UpdateReservation)
}

func (rvh *ReservationsHandler) FetchAllReservations(c echo.Context) error {
	reservations, _ := rvh.ReservationsUsecase.FetchAll(c.Request().Context())
	return c.JSON(http.StatusOK, reservations)
}

func (rvh *ReservationsHandler) FetchReservationByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	reservations, err := rvh.ReservationsUsecase.FetchByID(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "record not found" {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, reservations)
}

func (rvh *ReservationsHandler) CreateReservation(c echo.Context) error {
	reservation := domain.Reservations{}
	c.Bind(&reservation)

	err := rvh.ReservationsUsecase.Store(c.Request().Context(), &reservation)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success Create New Reservation",
	})
}

func (rvh *ReservationsHandler) UpdateReservation(c echo.Context) error {
	reservation := domain.Reservations{}
	c.Bind(&reservation)
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := rvh.ReservationsUsecase.FetchByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Data not found",
		})
	}

	err = rvh.ReservationsUsecase.Update(c.Request().Context(), &reservation, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success Update Reservation",
	})
}
