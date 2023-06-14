package api

import (
	"github.com/gin-gonic/gin"
	"senao-auth-srv/util"
)

type Server struct {
	config util.Config
	router *gin.Engine
}

func New(config util.Config) (*Server, error) {
	srv := &Server{
		config: config,
	}

	srv.setupRouter()
	return srv, nil
}

func (srv *Server) setupRouter() {
	router := gin.Default()

	router.POST("/register", srv.createAccount)
	router.POST("/verify", srv.verifyAccount)

	srv.router = router
}

func (srv *Server) Start(addr string) error {
	return srv.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
