package user

import (
	p "uller/src/permission"
	pl "uller/src/plan"
	r "uller/src/role"
)

type User struct {
  Name string `bson:"name" json:"name" form:"name"`
  Phone string `bson:"phone" json:"phone" form:"phone"`
  Email string `bson:"email" json:"email" form:"email"`
  Permissions []*p.Permission `bson:"permissions" json:"permissions"`
  Role *r.Role `bson:"role" json:"role"`
  Plan *pl.Plan `bson:"plan" json:"plan"`
  Password string `bson:"password" json:"password,omitempty" form:"password"`
}

