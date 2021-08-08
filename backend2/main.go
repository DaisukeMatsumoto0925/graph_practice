package main

import (
	"context"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	"github.com/DaisukeMatsumoto0925/backend2/src/graphql/resolver"
	"github.com/DaisukeMatsumoto0925/backend2/src/infra/rdb"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	db, err := rdb.InitDB()
	if err != nil {
		panic(err.Error())
	}

	middlewares := []echo.MiddlewareFunc{
		middleware.Recover(),
		middleware.Logger(),
		authorize(),
		NewCors(),
	}

	e.Use(middlewares...)

	graphqlHandler := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolver.Resolver{DB: db.Debug()},
			},
		),
	)

	playgroundHandler := playground.Handler("GraphQL", "/query")

	e.POST("/query", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/query", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.HideBanner = true
	e.Logger.Fatal(e.Start(":3000"))
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
