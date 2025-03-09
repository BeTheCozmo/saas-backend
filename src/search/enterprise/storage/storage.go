package enterprise_storage

import (
	"fmt"
	"strings"
	e "uller/src/search/enterprise"

	"gorm.io/gorm"
)

type EnterpriseStorage struct {
  DB              *gorm.DB
  legalEntity     string
  legalEntityBase string
  pageSize        uint64
}

func New() *EnterpriseStorage {
  return &EnterpriseStorage{}
}

func (es *EnterpriseStorage) Configure(db *gorm.DB) {
  es.DB = db
  es.legalEntity = "LegalEntity"
  es.legalEntityBase = "LegalEntityBase"
  es.pageSize = 10
}

func (es *EnterpriseStorage) GetByDocument(document string) *e.EnterpriseBase {
  var enterpriseBase e.EnterpriseBase
  es.DB.Table(es.legalEntityBase).Where("document = ?", document).Take(&enterpriseBase)
  return &enterpriseBase
}

func (es *EnterpriseStorage) GetEnterpriseHistory(document string, limit uint16) *[]e.Enterprise {
  var enterprises []e.Enterprise
  es.DB.Table(es.legalEntity).Where("document = ?", document).Order("createdAt DESC").Limit(int(limit)).Find(&enterprises)
  return &enterprises
}

func (es *EnterpriseStorage) GetEnterprisesByFilters(params *e.GetEnterprisesFilter) (*[]e.Enterprise, error) {
  var enterprises []e.Enterprise

  offset := 0
  if params.Page > 0 {
    offset = int((params.Page - 1) * es.pageSize)
  }

  subQuery := es.DB.Table(es.legalEntity).
    Select("DISTINCT document")

  // Query principal com JOIN correto
  query := es.DB.Table("(?) AS distinct_documents", subQuery).
    Joins("JOIN " + es.legalEntity + " ON " + es.legalEntity + ".document = distinct_documents.document").
    Select(es.legalEntity + ".*")

  // Apply filters
  for field, filter := range params.Filters {
    switch v := filter.(type) {
    case []string: // Agora esperamos []string diretamente
      if len(v) > 0 {
        query = query.Where(fmt.Sprintf("%s IN (?)", field), v)
      }
    case map[string]float64: // Agora esperamos map[string]float64 diretamente
      if min, ok := v["min"]; ok {
        query = query.Where(fmt.Sprintf("%s >= ?", field), min)
      }
      if max, ok := v["max"]; ok {
        query = query.Where(fmt.Sprintf("%s <= ?", field), max)
      }
    default:
      fmt.Printf("Não encontrado tipo do filtro para %s: %T\n", field, v)
    }
  }
  // Always exclude bankrupt entities
  query = query.Where("bankrupt = ?", false)

  // Apply ordering
  if len(params.Orders) > 0 {
    for field, direction := range params.Orders {
      if direction == "asc" || direction == "desc" {
        query = query.Order(fmt.Sprintf("%s %s", field, strings.ToUpper(direction)))
      }
    }
  }

  query = query.Order("id DESC")

  result := query.
    Limit(int(es.pageSize)).
    Offset(offset).
    Debug().
    Find(&enterprises)

  if result.Error != nil {
    return nil, result.Error
  }

  return &enterprises, nil
}

func (es *EnterpriseStorage) GetSuggestions(params *e.GetSuggestions) ([]string, error) {
  query := es.DB.Table(es.legalEntity).
    Select(fmt.Sprintf("DISTINCT `%s`", params.Field))

  // Aplicar filtros
  for field, filter := range params.Filters {
    if len(filter) > 0 {
      query = query.Where(fmt.Sprintf("`%s` IN (?)", field), filter)
    }
  }

  query = query.Where(fmt.Sprintf("`%s` LIKE ?", params.Field), params.Term+"%")

  var results []string
  result := query.
    Limit(10).
    Debug(). // Para ver a query executada
    Pluck(params.Field, &results)

  if result.Error != nil {
    return nil, result.Error
  }

  return results, nil
}
