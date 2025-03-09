package management_routes

import (
	m "uller/src/management"
	a "uller/src/user/authorization"

	"github.com/gin-gonic/gin"
)

type ManagementRoutes struct {
  router *gin.Engine
  m *m.Management
  a *a.Authorization
}

func New() *ManagementRoutes {
  return &ManagementRoutes{}
}

func (mr *ManagementRoutes) Configure(router *gin.Engine, m *m.Management, a *a.Authorization) {
  mr.router = router
  mr.m = m
  mr.a = a
}

func (mr *ManagementRoutes) ConfigureRoutes() {}