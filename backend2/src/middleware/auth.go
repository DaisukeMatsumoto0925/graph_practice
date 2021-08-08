package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo"
)

type key string

const (
	tokenKey key = "token"
)

func setToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

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
			newCtx := setToken(ctx.Request().Context(), tokenString)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return h(ctx)
		}
	}
}
