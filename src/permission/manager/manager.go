package permission_manager

import (
	"errors"
	"strings"
	u "uller/src/user"
	utils "uller/src/utils"

	setv2 "github.com/deckarep/golang-set/v2"
)

type PermissionManager struct {
}

func New() *PermissionManager {
  return &PermissionManager{}
}

func (pm *PermissionManager) Configure() {}

func GetAllGrantedPermissions(user *u.User) ([]string, error) {
  if user == nil {
    return nil, errors.New("erro ao buscar permissões: usuário não existe")
  }

  set := setv2.NewSet("")

  for _, v := range user.Permissions {
    set.Add(v.Name)
  }
  for _, v := range user.Role.Permissions {
    set.Add(v.Name)
  }
  for _, v := range user.Plan.Permissions {
    set.Add(v.Name)
  }

  return set.ToSlice(), nil
}

func GetAllGrantedPermissionsAndNotExpired(user *u.User) ([]string, error) {
  if user == nil {
    return nil, errors.New("erro ao buscar permissões: usuário não existe")
  }

  set := setv2.NewSet("")

  for _, v := range user.Permissions {
    if v.IsLimited && v.RemainingUses == 0 {
      continue;
    }
    set.Add(v.Name)
  }
  for _, v := range user.Role.Permissions {
    if v.IsLimited && v.RemainingUses == 0 {
      continue;
    }
    set.Add(v.Name)
  }
  for _, v := range user.Plan.Permissions {
    if v.IsLimited && v.RemainingUses == 0 {
      continue;
    }
    set.Add(v.Name)
  }

  return set.ToSlice(), nil
}

func HavePermission(user *u.User, permission string) bool {
  if user == nil {
    return false
  }
  permissions, err := GetAllGrantedPermissionsAndNotExpired(user);
  if err != nil {
    return false
  }
  return utils.Contains(permissions, permission)
}

func CanLogin(user *u.User) bool {
  return HavePermission(user, "user_login")
}

func GetAllEnterprisePermission(user *u.User) ([]string, error) {
  if user == nil {
    return nil, errors.New("user not provided")
  }

  set := setv2.NewSet("")

  for _, v := range user.Permissions {
    if v.IsLimited && v.RemainingUses == 0 {
      continue;
    }
    if strings.Contains(v.Name, "ent_") {
      set.Add(v.Name)
    }
  }
  for _, v := range user.Role.Permissions {
    if v.IsLimited && v.RemainingUses == 0 {
      continue;
    }
    if strings.Contains(v.Name, "ent_") {
      set.Add(v.Name)
    }
  }
  for _, v := range user.Plan.Permissions {
    if v.IsLimited && v.RemainingUses == 0 {
      continue;
    }
    if strings.Contains(v.Name, "ent_") {
      set.Add(v.Name)
    }
  }

  return set.ToSlice(), nil
}

func FilterEnterpriseSearchPermissions(permissions []string) ([]string, error) {
  if permissions == nil {
    return nil, errors.New("permissions not provided")
  }

  set := setv2.NewSet("")

  for _, v := range permissions {
    if strings.Contains(v, "ent_search_") {
      field := strings.TrimPrefix(v, "ent_search_")
      set.Add(field)
    }
  }
  
  return set.ToSlice(), nil
}

func FilterEnterpriseGetPermissions(permissions []string) ([]string, error) {
  if permissions == nil {
    return nil, errors.New("permissions not provided")
  }

  set := setv2.NewSet("")

  for _, v := range permissions {
    if strings.HasPrefix(v, "ent_get_") {
      field := strings.TrimPrefix(v, "ent_get_")
      set.Add(field)
    }
  }
  
  return set.ToSlice(), nil
}

func FilterEnterpriseFields(permissions []string) ([]string, error) {
  if permissions == nil {
    return nil, errors.New("permissions not provided")
  }

  set := setv2.NewSet("")

  for _, v := range permissions {
    splitedPermission := strings.Split(v, "_")
    length := len(splitedPermission)
    field := splitedPermission[length-1]
    set.Add(field)
  }

  return set.ToSlice(), nil
}