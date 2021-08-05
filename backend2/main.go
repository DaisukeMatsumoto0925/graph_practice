package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	"github.com/DaisukeMatsumoto0925/backend2/src/graphql/resolver"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	db, err := config.InitDB()
	if err != nil {
		panic(err.Error())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("CORS_ALLOW_ORIGIN")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(authorize())
	e.Use(injectStoreStatusLoader(db.Debug()))

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	hasRole := func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (interface{}, error) {
		token := getToken(ctx)
		if *token != role.String() {
			return nil, fmt.Errorf("Access denied")
		}
		fmt.Println("authenticate here !")
		return next(ctx)
	}
	graphqlHandler := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &resolver.Resolver{DB: db.Debug()},
				Directives: generated.DirectiveRoot{
					HasRole: hasRole,
				},
			},
		),
	)

	playgroundHandler := playground.Handler("GraphQL", "/query")

	e.POST("/query", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.Logger.SetLevel(elog.INFO)
	e.HideBanner = true
	e.Logger.Fatal(e.Start(":3000"))
}

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

func injectStoreStatusLoader(db *gorm.DB) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			loader := dataloader.CreateUserLoader(db)
			newCtx := dataloader.SetUserLoader(ctx.Request().Context(), loader)

			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return h(ctx)
		}
	}
}
