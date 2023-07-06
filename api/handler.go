package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"

	"senao-auth-srv/docs"
	"senao-auth-srv/middleware"
	"senao-auth-srv/service"
	"senao-auth-srv/util"
)

type Handler struct {
	config         util.Config
	accountService service.AccountService
}

var Module = fx.Options(
	fx.Provide(
		NewHandler,
	),
	fx.Invoke(registerService),
)

func NewHandler(config util.Config, accountService service.AccountService) *Handler {
	handler := &Handler{
		config:         config,
		accountService: accountService,
	}
	return handler
}

func registerService(router *gin.Engine, h *Handler) {
	router.Use(middleware.ErrorHandler())

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		v1.POST("/register", h.createAccount)
		v1.POST("/verify", h.verifyAccount)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
