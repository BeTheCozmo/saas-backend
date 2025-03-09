package server

import (
	er "uller/src/server/enterprise"
	ur "uller/src/server/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
  engine *gin.Engine
  ur *ur.UserRoutes
  er *er.EnterpriseRoutes
}

func New() *Server {
  return &Server{}
}

func (s *Server) Configure(router *gin.Engine, ur *ur.UserRoutes, er *er.EnterpriseRoutes) {
  s.engine = router
  s.ur = ur
  s.er = er
  s.configureRoutes()
}

func (s *Server) configureRoutes() {
  s.engine.Use(cors.Default())
  s.ur.ConfigureRoutes()
  s.er.ConfigureRoutes()
}

func (s *Server) Start(port string) {
  s.engine.Run(port)
}