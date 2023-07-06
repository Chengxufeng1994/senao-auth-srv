package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"senao-auth-srv/api"
	_db "senao-auth-srv/db"
	_redis "senao-auth-srv/redis"
	"senao-auth-srv/repository"
	"senao-auth-srv/server"
	"senao-auth-srv/service"
	"senao-auth-srv/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	app := fx.New(
		fx.Supply(config),
		_redis.Module,
		_db.Module,
		server.Module,
		api.Module,
		service.Module,
		repository.AccountRepoModule,
		fx.Invoke(startHttpServer),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal().Msgf("main: cannot start app err: %s", err)
		return
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-stopChan
	log.Info().Msg("main: shutting down server...")

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Fatal().Msgf("main: cannot stop app err: %s", err)
	}
}

func startHttpServer(lc fx.Lifecycle, config util.Config, gin *gin.Engine) *http.Server {
	addr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: gin,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msg("starting HTTP server")
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("stopping HTTP server")
			return nil
		},
	})

	return srv
}
