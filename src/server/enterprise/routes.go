package enterprise_routes

import (
	"fmt"
	"uller/src/logger"
	sm "uller/src/search/manager"
	a "uller/src/user/authorization"
	um "uller/src/user/manager"

	"github.com/gin-gonic/gin"
)

type EnterpriseRoutes struct {
  logger *logger.Logger
  sm *sm.SearchManager
  router *gin.Engine
  a *a.Authorization
  um *um.UserManager
}

func New() *EnterpriseRoutes {
  return &EnterpriseRoutes{}
}

func (er *EnterpriseRoutes) Configure(router *gin.Engine, sm *sm.SearchManager, a *a.Authorization, logger *logger.Logger, um *um.UserManager) {
  er.logger = logger
  er.router = router
  er.sm = sm
  er.a = a
  er.um = um
}

func (er *EnterpriseRoutes) ConfigureRoutes() {
  const prefix string = "enterprise"
  er.router.POST(fmt.Sprintf("%v/", prefix), er.getEnterprises)
  er.router.POST(fmt.Sprintf("%v/suggestion", prefix), er.getFilterSuggestions)
  er.router.GET(fmt.Sprintf("%v/:document", prefix), er.getByDocument)
  er.router.GET(fmt.Sprintf("%v/:document/history", prefix), er.getHistory)
}