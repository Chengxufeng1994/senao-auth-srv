package db

import (
	"github.com/go-redis/redis"
)

type Database struct {
	Client *redis.Client
}

func New(address string, password string) *Database {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return &Database{
		Client: client,
	}
}

func (db *Database) Conn() error {
	_, err := db.Client.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}
