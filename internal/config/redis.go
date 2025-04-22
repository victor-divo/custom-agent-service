package config

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	return rdb
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
