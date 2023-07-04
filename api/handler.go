package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"senao-auth-srv/db"
	"senao-auth-srv/docs"
	"senao-auth-srv/middleware"
	"senao-auth-srv/repository"
	"senao-auth-srv/service"
	"senao-auth-srv/util"
)

type Handler struct {
	config util.Config
	store  *db.Database
	Router *gin.Engine
}

func NewHandler(config util.Config, store *db.Database) *Handler {
	handler := &Handler{
		config: config,
		store:  store,
	}
	handler.setRoutes()
	return handler
}

func (h *Handler) setRoutes() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	accountRepo := repository.NewAccountRepoImpl(h.store)
	accountSvc := service.NewAccountServiceImpl(accountRepo)

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		NewAccountHandler(v1, accountSvc)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	h.Router = router
}
