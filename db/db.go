package db

import (
	"github.com/go-redis/redis"
)

type Database struct {
	Client *redis.Client
}

func New(address string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Database{
		Client: client,
	}, nil
}
