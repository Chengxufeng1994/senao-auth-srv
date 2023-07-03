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

// @BasePath /api/v1

type Server struct {
	config   util.Config
	database *db.Database
	router   *gin.Engine
}

func New(config util.Config, database *db.Database) (*Server, error) {
	srv := &Server{
		config:   config,
		database: database,
	}

	srv.setupRouter()
	return srv, nil
}

func (srv *Server) setupRouter() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	accountRepo := repository.NewAccountRepoImpl(srv.database)
	accountSvc := service.NewAccountServiceImpl(accountRepo)

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := router.Group("/api/v1")
	{
		NewAccountHandler(v1, accountSvc)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	srv.router = router
}

func (srv *Server) Start(addr string) error {
	return srv.router.Run(addr)
}
