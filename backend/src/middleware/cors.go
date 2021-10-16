package middleware

import (
	"github.com/labstack/echo"
)

func NewCors() echo.MiddlewareFunc {
	// return middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{os.Getenv("CORS_ALLOW_ORIGIN")},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// })

	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Writer.Header().Set("Access-Control-Allow-Origin", c.Request().Header.Get("Origin"))
			c.Response().Header().Set("Access-Control-Max-Age", "12h0m0s")
			c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Authorization")
			c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			return h(c)
		}
	}
}
