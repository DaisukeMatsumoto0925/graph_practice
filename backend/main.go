package main

import (
	"context"

	"github.com/DaisukeMatsumoto0925/backend/src/graphql/directive"
	"github.com/DaisukeMatsumoto0925/backend/src/graphql/resolver"
	"github.com/DaisukeMatsumoto0925/backend/src/graphql/subscriber"
	"github.com/DaisukeMatsumoto0925/backend/src/infra/rdb"
	"github.com/DaisukeMatsumoto0925/backend/src/infra/redis"
	"github.com/DaisukeMatsumoto0925/backend/src/infra/server"
	"github.com/DaisukeMatsumoto0925/backend/src/middleware"
	"github.com/labstack/echo"
)

func main() {
	db, err := rdb.InitDB()
	if err != nil {
		panic(err.Error())
	}
	redis := redis.New()

	subscribers := resolver.Subscribers{
		Message:    subscriber.NewMessageSubscriber(context.Background(), redis),
		UserStatus: subscriber.NewUserStatusSubscriber(context.Background(), redis),
	}
	loader := middleware.NewDataloader(db)

	middlewares := []echo.MiddlewareFunc{
		middleware.Authorize(),
		middleware.NewCors(),
		loader.InjectStoreStatusLoader(),
	}

	resolver := resolver.New(db, subscribers)
	directive := directive.New(db)
	graphqlHandler := server.GraphqlHandler(resolver, directive)
	router := server.NewRouter(graphqlHandler, middlewares)
	server.Run(router)

}
