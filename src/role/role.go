package role

import (
	p "uller/src/permission"
)

type Role struct {
  Name string `json:"name" bson:"name"`
  Permissions []*p.Permission `json:"permissions" bson:"permissions"`
}