package user_routes

import (
	"fmt"
	a "uller/src/user/authorization"
	um "uller/src/user/manager"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
  um *um.UserManager
  router *gin.Engine
  a *a.Authorization
}

func New() *UserRoutes {
  return &UserRoutes{}
}

func (ur *UserRoutes) Configure(router *gin.Engine, um *um.UserManager, a *a.Authorization) {
  ur.um = um
  ur.router = router
  ur.a = a
}

func (ur *UserRoutes) ConfigureRoutes() {
  const prefix string = "user"
  ur.router.POST(fmt.Sprintf("%v", prefix), ur.createUserHandler)
  ur.router.POST(fmt.Sprintf("%v/login", prefix), ur.authenticateHandler)
  ur.router.POST(fmt.Sprintf("%v/validate", prefix), ur.validateUser)
  ur.router.GET(fmt.Sprintf("%v/:credential", prefix), ur.getUserDetails)
}