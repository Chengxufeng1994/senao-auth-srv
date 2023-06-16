package main

import (
	"fmt"
	"os"
	"time"

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

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	redisAddr := fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort)
	redisPassword := config.RedisPassword
	database := db.New(redisAddr, redisPassword)
	dbRetryConn := 3
	isDbConn := false
	for i := 0; i < dbRetryConn; i++ {
		err = database.Conn()
		if err != nil {
			isDbConn = false
			log.Error().Msgf("retry connect database %d", i+1)
		} else {
			isDbConn = true
			break
		}
		time.Sleep(5 * time.Second)
	}
	if !isDbConn {
		log.Fatal().Msg("cannot connect database")
	}

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
