package config

import (
	"MyApp/global"
	"log"
	"context"

	"github.com/go-redis/redis/v8"
)

func InitRedis(){

	addr := AppConfig.Redis.Addr
	db := AppConfig.Redis.DB
	password := AppConfig.Redis.Password 

	RedisClient := redis.NewClient(&redis.Options{
		Addr: addr,
		DB: db,
		Password: password,
	})

	_, err := RedisClient.Ping(context.Background()).Result()

	if err !=nil{
		log.Fatalf("Failed to connect to Redis, got error: %v", err)
	}

	global.RedisDB = RedisClient
}