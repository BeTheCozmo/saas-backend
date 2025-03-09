package plan

import (
	p "uller/src/permission"
)

type Plan struct {
  Name string `bson:"name" json:"name"`
  Permissions []*p.Permission `bson:"permissions" json:"permissions"`
  StartedAt uint64 `bson:"startedAt" json:"startedAt"`
  EndAtBase uint64 `bson:"endAtBase" json:"endAtBase"`
  PriceBase uint64 `bson:"priceBase" json:"priceBase"`
}