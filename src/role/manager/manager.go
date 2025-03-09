package role_manager

import (
	r "uller/src/role"
	rs "uller/src/role/storage"
)

type RoleManager struct {
  rs *rs.RoleStorage
}

func New() *RoleManager {
  return &RoleManager{}
}

func (rm *RoleManager) Configure(rs *rs.RoleStorage) {
  rm.rs = rs
}

func (rm *RoleManager) GetRoleByName(name string) *r.Role {
  return rm.rs.GetRoleByName(name)
}

func (rm *RoleManager) GetUserRole() *r.Role {
  return rm.rs.GetRoleByName("user")
}

func (rm *RoleManager) GetAdminRole() *r.Role {
  return rm.rs.GetRoleByName("admin")
}