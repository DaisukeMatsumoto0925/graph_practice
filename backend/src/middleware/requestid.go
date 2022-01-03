package middleware

import (
	"context"

	"github.com/DaisukeMatsumoto0925/backend/src/infra/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type key string

const (
	REQUEST_ID_KEY key = "requestID"
)

func NewInjectRequestID() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			newCtx := context.WithValue(c.Request().Context(), REQUEST_ID_KEY, requestID)
			c.SetRequest(c.Request().WithContext(newCtx))
			return h(c)
		}
	}
}

func GetRequestID(ctx context.Context) string {
	requestID := ctx.Value(REQUEST_ID_KEY)

	if target, ok := requestID.(string); ok {
		return target
	} else {
		logger.Fatal("cannot get requestID from Context")
		panic("")
	}
}
