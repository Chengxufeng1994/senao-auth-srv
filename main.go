package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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

	handler := api.NewHandler(config, database)
	runHttpSrv(config, handler)
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

func runHttpSrv(config util.Config, handler *api.Handler) {
	addr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Info().Msgf("cannot start server error: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("server shutdown error: %s\n", err)
	}
	log.Info().Msg("server exiting")
}
