package permission

type Permission struct {
  Name string `json:"name" bson:"name"`
  IsLimited bool `json:"isLimited" bson:"isLimited"`
  RemainingUses uint32 `json:"remainingUses" bson:"remainingUses"`
}