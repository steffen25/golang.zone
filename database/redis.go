package database

import (
	"github.com/go-redis/redis"
	"github.com/steffen25/golang.zone/config"
	"strconv"
)

type RedisDB struct {
	*redis.Client
}

// Use this in the app.go's New function and add the struct type in the App struct in app.go
func NewRedisDB(dbCfg config.RedisConfig) (*RedisDB, error) {
	port := strconv.Itoa(dbCfg.Post)
	client := redis.NewClient(&redis.Options{
		Addr:     dbCfg.Host+":"+port,
		Password: "", // no password set
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisDB{client}, err
}