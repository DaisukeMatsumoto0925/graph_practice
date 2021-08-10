package middleware

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewCors() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("CORS_ALLOW_ORIGIN")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})
}
