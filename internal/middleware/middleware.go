package middleware

import (
	"os"

	"github.com/Alterra-DataOn-Kelompok-5/reservation-service/internal/dto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LogMiddlewares(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           " ${time_custom} | ${host} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} \n",
		CustomTimeFormat: "2006/01/02 15:04:05",
		Output:           os.Stdout,
	}))
}

func JWTMiddleware(claims dto.JWTClaims, signingKey []byte) echo.MiddlewareFunc {
	config := middleware.JWTConfig{
		Claims:     &dto.JWTClaims{},
		SigningKey: signingKey,
	}
	return middleware.JWTWithConfig(config)
}
