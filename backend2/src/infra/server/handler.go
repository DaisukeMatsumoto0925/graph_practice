package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	"github.com/labstack/echo"
)

func GraphqlHandler(resolver generated.ResolverRoot) echo.HandlerFunc {
	c := generated.Config{
		Resolvers: resolver,
		// Directives: directive,
	}

	h := handler.New(
		generated.NewExecutableSchema(c),
	)
	h.AddTransport(transport.POST{}) // https://zenn.dev/konboi/articles/ee8ec5c27b98576de3db

	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func playgroundHandler() echo.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
