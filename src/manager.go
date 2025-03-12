package manager

import (
	"os"
	"uller/src/logger"
	m "uller/src/management"
	pm "uller/src/permission/manager"
	ps "uller/src/plan/storage"
	rm "uller/src/role/manager"
	rs "uller/src/role/storage"
	es "uller/src/search/enterprise/storage"
	sm "uller/src/search/manager"
	s "uller/src/server"
	er "uller/src/server/enterprise"
	ur "uller/src/server/user"
	a "uller/src/user/authorization"
	um "uller/src/user/manager"
	us "uller/src/user/storage"

	gormlogger "gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Manager struct {
  Logger *logger.Logger
  UserStorage *us.UserStorage
  UserManager *um.UserManager
  PermissionManager *pm.PermissionManager
  Authorization *a.Authorization
  Server *s.Server
  RoleStorage *rs.RoleStorage
  RoleManager *rm.RoleManager
  PlanStorage *ps.PlanStorage
  EnterpriseStorage *es.EnterpriseStorage
  EnterpriseRoutes *er.EnterpriseRoutes
  SearchManager *sm.SearchManager
  UserRoutes *ur.UserRoutes
  Management *m.Management
}

func New() *Manager {
  return &Manager{};
}
func (u *Manager) CreateModules() {
  u.Logger = logger.New()
  u.Authorization = a.New()
  u.PermissionManager = pm.New()
  u.PlanStorage = ps.New()
  u.RoleStorage = rs.New()
  u.RoleManager = rm.New()
  u.UserStorage = us.New()
  u.UserManager = um.New()
  u.UserRoutes = ur.New()
  u.Server = s.New()
  u.EnterpriseStorage = es.New()
  u.EnterpriseRoutes = er.New()
  u.SearchManager = sm.New()
  u.Management = m.New()
}
func (u *Manager) ConfigureModules() {
  err := godotenv.Load()
  u.Logger.Configure()
  if err != nil {
    u.Logger.Error("error loading .env")
  }
  
  enterpriseDSN := os.Getenv("ENTERPRISES_DB_URL")
  enterpriseClient, err := gorm.Open(mysql.Open(enterpriseDSN), &gorm.Config{
    Logger: gormlogger.Default.LogMode(gormlogger.Info),
  })
  if err != nil {
    u.Logger.Debug(err)
  }
  
  ullerDSN := os.Getenv("ULLER_DB_URL")
  if ullerDSN == "" {
    u.Logger.Debug("uller dsn is empty")
  }
  ullerDBName := os.Getenv("ULLER_DB_NAME")
  if ullerDBName == "" {
    ullerDBName = "uller"
  }

  ullerClient, err := mongo.Connect(options.Client().ApplyURI(ullerDSN))
  if err != nil {
    u.Logger.Debug(err)
  }
  
  jwtSecret := os.Getenv("JWT_SECRET")
  if jwtSecret == "" {
    jwtSecret = "53ArF`*U7/viV!xr"
  }
  router := gin.Default()

  u.Management.Configure()
  u.PermissionManager.Configure()
  u.EnterpriseStorage.Configure(enterpriseClient)
  u.SearchManager.Configure(u.PermissionManager, u.EnterpriseStorage, u.Logger, u.UserManager)
  u.EnterpriseRoutes.Configure(router, u.SearchManager, u.Authorization, u.Logger, u.UserManager)
  u.Authorization.Configure(jwtSecret)
  u.RoleStorage.Configure(ullerClient, ullerDBName)
  u.RoleManager.Configure(u.RoleStorage)
  u.PlanStorage.Configure(ullerClient, ullerDBName)
  u.UserStorage.Configure(ullerClient, ullerDBName, u.Logger)
  u.UserManager.Configure(u.UserStorage, u.PermissionManager, u.Authorization, u.RoleManager, u.PlanStorage, u.Logger)
  u.UserRoutes.Configure(router, u.UserManager, u.Authorization)
  u.Server.Configure(router, u.UserRoutes, u.EnterpriseRoutes)
}
func (u *Manager) Run() {
  u.Server.Start("0.0.0.0:8080")
}
