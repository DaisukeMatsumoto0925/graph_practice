package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewRouter(graphqlHandler echo.HandlerFunc, middlewares []echo.MiddlewareFunc) *echo.Echo {
	r := echo.New()
	r.Use(middleware.Recover())
	// r.Use(middleware.Logger())
	r.Use(middlewares...)

	r.GET("/", playgroundHandler())

	r.POST("/query", graphqlHandler)
	r.GET("/query", graphqlHandler)

	return r
}
