package middleware

import (
	"strings"

	"github.com/DaisukeMatsumoto0925/backend2/src/util/appcontext"
	"github.com/labstack/echo"
)

func Authorize() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeaderParts := strings.Split(ctx.Request().Header.Get("Authorization"), " ")
			if len(authHeaderParts) < 2 {
				return h(ctx)
			}
			tokenString := authHeaderParts[1]
			if tokenString == "" {
				return h(ctx)
			}
			newCtx := appcontext.SetToken(ctx.Request().Context(), tokenString)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return h(ctx)
		}
	}
}
