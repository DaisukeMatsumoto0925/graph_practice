package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Client = redis.Client

var Nil = redis.Nil

func New() *Client {
	dbConf, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbConf,
	})

	fmt.Println(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DB"))
	context := context.Background()

	err := client.Ping(context).Err()

	if err != nil {
		log.Fatal("failed to connect redis", err)
	}

	log.Print("success to connect redis!")

	return client
}
