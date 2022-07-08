package middleware

import (
	"os"

	echo "github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo) {
	e.Use(
		echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format:           " ${time_custom} | ${host} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} \n",
			CustomTimeFormat: "2006/01/02 15:04:05",
			Output:           os.Stdout,
		}),
	)
}
