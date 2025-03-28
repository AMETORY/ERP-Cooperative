package services

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var (
	REDIS = &redis.Client{}
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err)
	}

	REDIS = rdb
	return rdb
}
