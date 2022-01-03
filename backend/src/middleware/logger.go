package middleware

import (
	"github.com/DaisukeMatsumoto0925/backend/src/infra/logger"
	"github.com/DaisukeMatsumoto0925/backend/src/util/appcontext"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func NewInjectLogger() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := GetRequestID(c.Request().Context())
			logger := logger.New().WithFields(logrus.Fields{"requestId": requestID})
			newCtx := appcontext.SetLogger(c.Request().Context(), logger)
			c.SetRequest(c.Request().WithContext(newCtx))
			return h(c)
		}
	}
}
