package database

import "github.com/go-redis/redis"

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
