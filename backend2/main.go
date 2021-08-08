package main

import (
	"context"
	"os"
	"strings"

	"github.com/DaisukeMatsumoto0925/backend2/src/graphql/resolver"
	"github.com/DaisukeMatsumoto0925/backend2/src/infra/rdb"
	"github.com/DaisukeMatsumoto0925/backend2/src/infra/server"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	db, err := rdb.InitDB()
	if err != nil {
		panic(err.Error())
	}

	middlewares := []echo.MiddlewareFunc{
		authorize(),
		NewCors(),
	}

	resolver := resolver.New(db)
	graphqlHandler := server.GraphqlHandler(resolver)
	router := server.NewRouter(graphqlHandler, middlewares)
	server.Run(router)

}

// ---middleware and more------------------------------------------------------------------------

func authorize() echo.MiddlewareFunc {
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

func NewCors() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("CORS_ALLOW_ORIGIN")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})
}

type key string

const (
	tokenKey key = "token"
)

func setToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func getToken(ctx context.Context) *string {
	token := ctx.Value(tokenKey)
	if token, ok := token.(string); ok {
		return &token
	}
	return nil
}
