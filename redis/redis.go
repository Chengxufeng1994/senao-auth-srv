package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"senao-auth-srv/util"
)

type Cache interface{}

type CacheImpl struct {
	Client *redis.Client
}

var Module = fx.Module(
	"redis",
	fx.Provide(
		NewRedisClient,
	),
	fx.Invoke(func(redis *CacheImpl) {}),
)

func NewRedisClient(lc fx.Lifecycle, config util.Config) *CacheImpl {
	addr := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.RedisPassword,
		DB:       0,
	})
	cache := &CacheImpl{
		Client: client,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := cache.Client.Ping().Result()
			if err != nil {
				log.Error().Msgf("connect to redis error: %v", err.Error())
				return err
			}
			log.Info().Msg("connect to redis successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("disconnect to redis")
			err := cache.Client.Close()
			if err != nil {
				log.Error().Msgf("disconnect to redis error: %v", err.Error())
				return err
			}
			log.Info().Msg("disconnect to redis successfully")
			return nil
		},
	})

	return cache
}
