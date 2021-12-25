package server

import (
	"github.com/DaisukeMatsumoto0925/backend/src/rest/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewRouter(graphqlHandler echo.HandlerFunc, middlewares []echo.MiddlewareFunc, ctrl *controller.Controller) *echo.Echo {
	r := echo.New()
	r.Use(middleware.Recover())
	// r.Use(middleware.Logger())
	r.Use(middlewares...)

	// GraphQLPlayGround
	r.GET("/", playgroundHandler())

	// GraphQLAPIs
	r.POST("/query", graphqlHandler)
	r.GET("/query", graphqlHandler)

	// RestAPIs
	apiV1 := r.Group("/api/v1")
	csvRouter(apiV1, ctrl)

	return r
}

func csvRouter(g *echo.Group, ctrl *controller.Controller) {
	constroller := controller.NewAnalyticsController(ctrl)
	g.GET("/analytics", func(e echo.Context) error {
		return constroller.TaskCSV(e)
	})
}
