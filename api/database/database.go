package database

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func CreateClient(dbNo int) *redis.Client {
	dbPort := os.Getenv("DB_ADDR")
	if dbPort == "" {
		dbPort = "db:6379"
	}
	fmt.Println(dbPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     dbPort,
		Password: "",
		DB:       dbNo,
	})
	return rdb
}
