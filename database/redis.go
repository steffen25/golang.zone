package database

import (
	"github.com/go-redis/redis"
	"github.com/steffen25/golang.zone/config"
	"strconv"
)

var client *redis.Client

type RedisDB struct {
	*redis.Client
}

func RedisConnection() (*redis.Client, error) {

	var err error

	if client == nil {
		client, err = NewClient()
	}

	return client, err
}

func NewClient() (*redis.Client, error) {

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()

	return client, err
}

// Use this in the app.go's New function and add the struct type in the App struct in app.go
func NewRedisDB(dbCfg config.RedisConfig) (*RedisDB, error) {
	port := strconv.Itoa(dbCfg.Post)
	client = redis.NewClient(&redis.Options{
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