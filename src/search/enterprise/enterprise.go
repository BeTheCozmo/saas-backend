package enterprise

type EnterpriseStatus string
const (
  ACTIVE EnterpriseStatus = "ACTIVE"
  SUSPENDED EnterpriseStatus = "SUSPENDED"
  UNAPPEALABLE EnterpriseStatus = "UNAPPEALABLE"
  CLOSED EnterpriseStatus = "CLOSED"
  NULL EnterpriseStatus = "NULL"
)

type EnterpriseBranchOffice string
const (
  BRANCH EnterpriseBranchOffice = "BRANCH"
  HEAD_OFFICE EnterpriseBranchOffice = "HEAD_OFFICE"
)

type EnterpriseSimples string
const (
  SIM EnterpriseSimples = "SIM"
  NAO EnterpriseSimples = "NAO"
  SEM_INFORMACAO EnterpriseSimples = "SEM_INFORMACAO"
)

type Enterprise struct {
  Id uint64 `gorm:"primaryKey, column:id" json:"id"`
  Document uint64 `gorm:"index, column:document" json:"document"`
  Status EnterpriseStatus `gorm:"column:status" json:"status"`
  StatusDate string `gorm:"column:statusDate" json:"statusDate"`
  SecondaryActivitiesIds string `gorm:"column:secondaryActivitiesIds" json:"secondaryActivitiesIds"`
  Mei bool `gorm:"column:mei" json:"mei"`
  BranchOffice EnterpriseBranchOffice `gorm:"column:branchOffice" json:"branchOffice"`
  Debt float64 `gorm:"column:debt" json:"debt"`
  // CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
  FantasyName string `gorm:"column:fantasyName" json:"fantasyName"`
  Name string `gorm:"column:name" json:"name"`
  Enriched bool `gorm:"column:enriched" json:"enriched"`
  Partners int16 `gorm:"column:partners" json:"partners"`
  Simples EnterpriseSimples `gorm:"column:simples" json:"simples"`
  State string `gorm:"column:state" json:"state"`
  City string `gorm:"column:city" json:"city"`
  District string `gorm:"column:district" json:"district"`
  Street string `gorm:"column:street" json:"street"`
  ZipCode string `gorm:"column:zipCode" json:"zipCode"`
  Email string `gorm:"column:email" json:"email"`
  FoundationDate string `gorm:"column:foundationDate" json:"foundationDate"`
  NatureCode uint32 `gorm:"column:natureCode" json:"natureCode"`
  NatureDescription string `gorm:"column:natureDescription" json:"natureDescription"`
  PrimaryActivityCode string `gorm:"column:primaryActivityCode" json:"primaryActivityCode"`
  PrimaryActivityDescription string `gorm:"column:primaryActivityDescription" json:"primaryActivityDescription"`
  SocialCapital string `gorm:"column:socialCapital" json:"socialCapital"`
  Telephone string `gorm:"column:telephone" json:"telephone"`
  LastDebt float64 `gorm:"column:lastDebt" json:"lastDebt"`
  Variation float64 `gorm:"column:variation" json:"variation"`
  VariationPercentage float64 `gorm:"column:variationPercentage" json:"variationPercentage"`
  CollectId uint64 `gorm:"column:collectId" json:"collectId"`
  Bankrupt bool `gorm:"column:bankrupt" json:"bankrupt"`
  PartnersBoardId int64 `gorm:"column:partnersBoardId" json:"partnersBoardId"`
}

type EnterpriseBase struct {
  Document uint64 `gorm:"primaryKey, column:document" json:"document"`
  Status EnterpriseStatus `gorm:"column:status" json:"status"`
  StatusDate string `gorm:"column:statusDate" json:"statusDate"`
  SecondaryActivitiesIds string `gorm:"column:secondaryActivitiesIds" json:"secondaryActivitiesIds"`
  Mei bool `gorm:"column:mei" json:"mei"`
  BranchOffice EnterpriseBranchOffice `gorm:"column:branchOffice" json:"branchOffice"`
  // CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
  FantasyName string `gorm:"column:fantasyName" json:"fantasyName"`
  Name string `gorm:"column:name" json:"name"`
  Enriched bool `gorm:"column:enriched" json:"enriched"`
  Partners int16 `gorm:"column:partners" json:"partners"`
  Simples EnterpriseSimples `gorm:"column:simples" json:"simples"`
  State string `gorm:"column:state" json:"state"`
  City string `gorm:"column:city" json:"city"`
  District string `gorm:"column:district" json:"district"`
  Street string `gorm:"column:street" json:"street"`
  ZipCode string `gorm:"column:zipCode" json:"zipCode"`
  Email string `gorm:"column:email" json:"email"`
  FoundationDate string `gorm:"column:foundationDate" json:"foundationDate"`
  NatureCode uint32 `gorm:"column:natureCode" json:"natureCode"`
  NatureDescription string `gorm:"column:natureDescription" json:"natureDescription"`
  PrimaryActivityCode string `gorm:"column:primaryActivityCode" json:"primaryActivityCode"`
  PrimaryActivityDescription string `gorm:"column:primaryActivityDescription" json:"primaryActivityDescription"`
  SocialCapital string `gorm:"column:socialCapital" json:"socialCapital"`
  Telephone string `gorm:"column:telephone" json:"telephone"`
  Bankrupt bool `gorm:"column:bankrupt" json:"bankrupt"`
  PartnersBoardId int64 `gorm:"column:partnersBoardId" json:"partnersBoardId"`
}