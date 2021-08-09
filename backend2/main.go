package main

import (
	"github.com/DaisukeMatsumoto0925/backend/src/graphql/directive"
	"github.com/DaisukeMatsumoto0925/backend/src/graphql/resolver"
	"github.com/DaisukeMatsumoto0925/backend/src/infra/rdb"
	"github.com/DaisukeMatsumoto0925/backend/src/infra/server"
	"github.com/DaisukeMatsumoto0925/backend/src/middleware"
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
