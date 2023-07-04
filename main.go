package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"senao-auth-srv/api"
	"senao-auth-srv/db"

	"senao-auth-srv/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if config.Environment == "development" {
		gin.SetMode(gin.DebugMode)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	redisAddr := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)
	redisPassword := config.RedisPassword
	database := db.New(redisAddr, redisPassword)

	const maxRetries int = 3
	const timeout = 3
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		var retries int
		for {
			_, err = database.Conn()
			if err == nil {
				log.Info().Msg("connect to database successfully")
				break
			}
			retries++
			if retries == maxRetries {
				log.Fatal().Msg("cannot connect to database")
			}
			log.Error().Msgf("connect to database - retrying... %d", retries)
			time.Sleep(timeout * time.Second)
		}
		wg.Done()
	}()
	wg.Wait()

	runGinSrv(config, database)
}

func runGinSrv(config util.Config, database *db.Database) {
	serverAddr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	srv, err := api.New(config, database)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = srv.Start(serverAddr)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
