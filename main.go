package main

import (
	"fmt"
	"os"
	"senao-auth-srv/api"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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

	runGinSrv(config)
}

func runGinSrv(config util.Config) {
	serverAddr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	srv, err := api.New(config)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = srv.Start(serverAddr)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
