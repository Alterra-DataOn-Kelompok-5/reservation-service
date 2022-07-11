package http

import (
	reservationStatus "github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/app/reservation_status"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/app/reservations"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/factory"
	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/pkg/util"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func NewHttp(e *echo.Echo, f *factory.Factory) {
	e.Validator = &util.CustomValidator{Validator: validator.New()}

	e.GET("/status", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})
	v1 := e.Group("/api/v1")
	reservations.NewHandler(f).Route(v1.Group("/reservations"))
	reservationStatus.NewHandler(f).Route(v1.Group("/status"))
}
