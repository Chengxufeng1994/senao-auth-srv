package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	"senao-auth-srv/util"
)

var Module = fx.Provide(
	NewGinServer,
)

func NewGinServer(config util.Config) *gin.Engine {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if config.Environment == "development" {
		gin.SetMode(gin.DebugMode)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return gin.Default()
}
