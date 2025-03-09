package user_manager

import (
	"errors"
	"fmt"
	"strings"
	"uller/src/logger"
	p "uller/src/permission"
	pm "uller/src/permission/manager"
	ps "uller/src/plan/storage"
	r "uller/src/role"
	rm "uller/src/role/manager"
	"uller/src/user"
	a "uller/src/user/authorization"
	us "uller/src/user/storage"
	"unicode"

	"uller/src/utils"
)

type UserManager struct {
  logger *logger.Logger
  us *us.UserStorage
  pm *pm.PermissionManager
  a *a.Authorization
  rm *rm.RoleManager
  ps *ps.PlanStorage
}

func New() *UserManager {
  return &UserManager{}
}

func (um *UserManager) Configure(
  us *us.UserStorage,
  pm *pm.PermissionManager,
  a *a.Authorization,
  rm *rm.RoleManager,
  ps *ps.PlanStorage,
  logger *logger.Logger) {
  
  um.logger = logger
  um.us = us
  um.pm = pm
  um.a = a
  um.rm = rm
  um.ps = ps

  um.upsertAdmins()
}

func (um *UserManager) upsertAdmins() {
  antonyEmail := "antony.cozmo@gmail.com"
  antony, err := um.us.GetByEmail(antonyEmail)
  if err != nil {
    um.logger.Error(err, antonyEmail)
  }
  if antony != nil {
    um.makeAdmin(antony)
  }
}

func (um *UserManager) makeAdmin(user *user.User) {
  adminRole := um.rm.GetAdminRole()
  um.changeUserRoleTo(user, adminRole)
}

func (um *UserManager) changeUserRoleTo(user *user.User, role *r.Role) {
  um.us.ChangeUserRoleTo(user, role)
}

func isEmail(credential string) bool {
  return strings.Contains(credential, "@") && strings.Contains(credential, ".")
}

func isNumeric(s string) bool {
  if s == "" {
    return false
  }

  for _, r := range s {
    if !unicode.IsDigit(r) {
      return false
    }
  }

  return true
}

func isPhone(credential string) bool {
  return isNumeric(credential)
}

func (um *UserManager) GetByCredential(credential string) (*user.User, error) {
  var user *user.User
  if isEmail(credential) {
    _user, err := um.us.GetByEmail(credential)
    if err != nil {
      return nil, err
    }
    user = _user
  } else if isPhone(credential) {
    _user, err := um.us.GetByPhone(credential)
    if err != nil {
      return nil, err
    }
    user = _user
  } else {
    return nil, errors.New("invalid e-mail or phone number")
  }

  if user == nil {
    return nil, errors.New("no user found with given credential")
  }

  return user, nil
}

func (um *UserManager) CreateUser(userData *user.User) (string, error) {
  if (userData.Email == "") && (userData.Phone == "") {
    return "", errors.New("email or phone don't provided")
  }
  hashedPassword, err := utils.HashPassword(userData.Password)
  if err != nil {
    return "", err
  }
  userData.Password = hashedPassword

  userPermissions := um.getBasePermissions()
  userRole := um.rm.GetUserRole()
  userPlan := um.ps.GetFreePlan()
  userData.Permissions = userPermissions
  userData.Role = userRole
  userData.Plan = userPlan

  err = um.us.Create(userData)
  if err != nil {
    return "", err
  }

  createdUser, err := um.us.GetByEmail(userData.Email)
  if err != nil {
    return "", err
  }

  token, err := um.a.GetAuthorizationToken(createdUser)
  if err != nil {
    fmt.Println("Error with authorization token", err)
    return "", err
  }

  return token, nil
}

func (um *UserManager) GetAuthorizationToken(user *user.User) (string, error) {
  return um.a.GetAuthorizationToken(user)
}

func (um *UserManager) GetAuthorization(credential string, password string) (string, error) {
  user, err := um.GetByCredential(credential)
  if err != nil {
    return "", err
  }

  if !utils.CheckPasswordHash(password, user.Password) {
    return "", errors.New("wrong password")
  }

  canLogin := pm.CanLogin(user)
  if !canLogin {
    return "", errors.New("you don't have permission to log-in")
  }

  token, err := um.GetAuthorizationToken(user)
  if err != nil {
    return "", err
  }

  err = um.UsePermissions(user, []string{"user_login"})
  if err != nil {
    return "", err
  }

  return token, nil
}

func (um *UserManager) getBasePermissions() []*p.Permission {
  return []*p.Permission{
    {
      Name: "user_login",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "user_get",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "user_change_name",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "user_change_profile",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "user_req_change_email",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "user_change_email",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "user_req_change_phone",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "user_change_phone",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "plan_get",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "plan_buy",
      IsLimited: false,
      RemainingUses: 0,
    },
    {
      Name: "plan_cancel",
      IsLimited: false,
      RemainingUses: 0,
    },
  }
}

func (um *UserManager) GetDetails(requester *user.User, credential string) (*user.User, error) {
  if !pm.HavePermission(requester, "user_get") {
    return nil, errors.New("you don't have permission to see details of other users")
  }

  user, err := um.GetByCredential(credential)
  if err != nil {
    return nil, err
  }

  user.Password = ""

  return user, nil
}

func (um *UserManager) UsePermissions(user *user.User, usedPermissions []string) error {
  for _, usedPermission := range usedPermissions {
    err := um.usePermission(user, usedPermission)
    if err != nil {
      return err
    }
  }

  return um.SaveUser(user)
}

func (um *UserManager) usePermission(user *user.User, usedPermission string) error {
  used, err := um.usePermissionForPermissions(usedPermission, &user.Plan.Permissions)
  if err != nil {
    return err
  } else if used {
    return nil
  }
  used, err = um.usePermissionForPermissions(usedPermission, &user.Role.Permissions)
  if err != nil {
    return err
  } else if used {
    return nil
  }
  used, err = um.usePermissionForPermissions(usedPermission, &user.Permissions)
  if err != nil {
    return err
  } else if used {
    return nil
  }
  
  return nil
}

func (um *UserManager) usePermissionForPermissions(usedPermission string, permissions *[]*p.Permission) (bool, error) {
  for _, permission := range *permissions {
    if permission.Name == usedPermission {
      if !permission.IsLimited {
        return false, nil
      }
      if permission.RemainingUses == 0 {
        errorMessage := fmt.Sprintf("you dont have permission to use '%v' permission", permission)
        return false, errors.New(errorMessage)
      }
      permission.RemainingUses = permission.RemainingUses - 1
      return true, nil
    }
  }
  return false, nil
}

func (um *UserManager) SaveUser(user *user.User) error {
  return um.us.SaveUser(user)
}