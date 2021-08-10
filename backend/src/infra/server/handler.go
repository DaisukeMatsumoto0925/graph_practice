package server

import (
	// "github.com/99designs/gqlgen/graphql/handler"
	// "github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/99designs/gqlgen/handler"
	"github.com/DaisukeMatsumoto0925/backend/graph/generated"
	"github.com/labstack/echo"
)

func GraphqlHandler(resolver generated.ResolverRoot, directive generated.DirectiveRoot) echo.HandlerFunc {
	c := generated.Config{
		Resolvers:  resolver,
		Directives: directive,
	}

	h := handler.GraphQL(
		generated.NewExecutableSchema(c),
	)

	// NOTE: handler.Newが推奨だが playgroundのdocsが読み込まれない？
	// h := handler.New(
	// h.AddTransport(transport.POST{}) // https://zenn.dev/konboi/articles/ee8ec5c27b98576de3db

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
