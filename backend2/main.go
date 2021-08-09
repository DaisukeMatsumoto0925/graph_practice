package main

import (
	"github.com/DaisukeMatsumoto0925/backend2/src/graphql/directive"
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

	loader := middleware.NewDataloader(db)

	middlewares := []echo.MiddlewareFunc{
		middleware.Authorize(),
		middleware.NewCors(),
		loader.InjectStoreStatusLoader(),
	}

	resolver := resolver.New(db)
	directive := directive.New(db)
	graphqlHandler := server.GraphqlHandler(resolver, directive)
	router := server.NewRouter(graphqlHandler, middlewares)
	server.Run(router)

}
