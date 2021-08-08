package main

import (
	"github.com/DaisukeMatsumoto0925/backend2/src/graphql/resolver"
	"github.com/DaisukeMatsumoto0925/backend2/src/infra/rdb"
	"github.com/DaisukeMatsumoto0925/backend2/src/infra/server"
	"github.com/DaisukeMatsumoto0925/backend2/src/middleware"
	"github.com/labstack/echo"
)

func main() {
	db, err := rdb.InitDB()
	if err != nil {
		panic(err.Error())
	}

	middlewares := []echo.MiddlewareFunc{
		middleware.Authorize(),
		middleware.NewCors(),
	}

	resolver := resolver.New(db)
	graphqlHandler := server.GraphqlHandler(resolver)
	router := server.NewRouter(graphqlHandler, middlewares)
	server.Run(router)

}

// ---middleware and more------------------------------------------------------------------------

// hasRole := func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (interface{}, error) {
// 	token := getToken(ctx)
// 	if *token != role.String() {
// 		return nil, fmt.Errorf("Access denied")
// 	}
// 	fmt.Println("authenticate here !")
// 	return next(ctx)
// }

// graphqlHandler := handler.NewDefaultServer(
// 	generated.NewExecutableSchema(
// 		generated.Config{
// 			Resolvers: &resolver.Resolver{DB: db.Debug()},
// 			Directives: generated.DirectiveRoot{
// 				HasRole: hasRole,
// 			},
// 		},
// 	),
// )

// func getToken(ctx context.Context) *string {
// 	token := ctx.Value(tokenKey)
// 	if token, ok := token.(string); ok {
// 		return &token
// 	}
// 	return nil
// }
