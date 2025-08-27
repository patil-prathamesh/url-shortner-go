package database

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("PORT"),
		Password: os.Getenv("PASSWORD"),
		DB:       dbNo,
	})
	return rdb
}
