package search_manager

import (
	"errors"
	"fmt"
	"uller/src/logger"
	pm "uller/src/permission/manager"
	e "uller/src/search/enterprise"
	es "uller/src/search/enterprise/storage"
	u "uller/src/user"
	um "uller/src/user/manager"
	"uller/src/utils"
)

type SearchManager struct {
  logger *logger.Logger
  es *es.EnterpriseStorage
  pm *pm.PermissionManager
  um *um.UserManager
}

func New() *SearchManager {
  return &SearchManager{}
}

func (sm *SearchManager) Configure(pm *pm.PermissionManager, es *es.EnterpriseStorage, logger *logger.Logger, um *um.UserManager) {
  sm.logger = logger
  sm.pm = pm
  sm.es = es
  sm.um = um
}

func (sm *SearchManager) GetEnterpriseBaseByDocument(document string, u *u.User) (*e.EnterpriseBase, error) {
  if !pm.HavePermission(u, "ent_details") {
    return nil, errors.New("you do not have permission to access the details of enterprises")
  }
  return sm.es.GetByDocument(document)
}

func (sm *SearchManager) GetEnterpriseBaseByDocumentMap(document string, u *u.User) (map[string]interface{}, error) {
  enterprise, err := sm.GetEnterpriseBaseByDocument(document, u)
  if err != nil {
    return nil, err
  }
  if enterprise == nil {
    return nil, fmt.Errorf("enterprise doesn't exists")
  }

  permissions, err := pm.GetAllGrantedPermissionsAndNotExpired(u)
  if err != nil {
    return nil, err
  }

  getPermissions, err := pm.FilterEnterpriseGetPermissions(permissions)
  if err != nil {
    return nil, err
  }

  enterpriseMap, err := utils.StructToMap(enterprise)
  if err != nil {
    return nil, err
  }

  usedPermissions := []string{"ent_details"}
  for key := range enterpriseMap {
    if !utils.Contains(getPermissions, key) {
      delete(enterpriseMap, key)
    }
    usedPermissions = append(usedPermissions, "ent_get_"+key)
  }
  err = sm.um.UsePermissions(u, usedPermissions)
  if err != nil {
    return nil, err
  }

  return enterpriseMap, nil
}

func (sm *SearchManager) GetEnterpriseHistoryMap(document string, collectInterval uint16, u *u.User) (*[]map[string]interface{}, error) {
  if !pm.HavePermission(u, "ent_history") {
    return nil, errors.New("you do not have permission to access the history of enterprises")
  }
  enterprises := sm.es.GetEnterpriseHistory(document, collectInterval)

  permissions, err := pm.GetAllGrantedPermissionsAndNotExpired(u)
  if err != nil {
    return nil, err
  }

  getPermissions, err := pm.FilterEnterpriseGetPermissions(permissions)
  if err != nil {
    return nil, err
  }

  var rEnterprises []map[string]interface{}
  for _, enterprise := range *enterprises {
    enterpriseMap, err := utils.StructToMap(enterprise)
    if err != nil {
      return nil, err
    }

    for key := range enterpriseMap {
      if !utils.Contains(getPermissions, key) {
        delete(enterpriseMap, key)
      }
    }

    rEnterprises = append(rEnterprises, enterpriseMap)
  }

  err = sm.um.UsePermissions(u, []string{"ent_history"})
  if err != nil {
    return nil, err
  }

  return &rEnterprises, nil
}

func (sm *SearchManager) GetEnterprisesByFilters(filters *e.GetEnterprisesFilter, u *u.User) (*[]map[string]interface{}, error) {
  permissions, err := pm.GetAllGrantedPermissionsAndNotExpired(u)
  if err != nil {
    sm.logger.Debug("Error getting permissions")
    return nil, err
  }
  sm.logger.Debug(permissions)

  searchPermissions, err := pm.FilterEnterpriseSearchPermissions(permissions)
  if err != nil {
    sm.logger.Debug("Error getting search permissions")
    return nil, err
  }
  sm.logger.Debug("search permissions:", searchPermissions)

  for permission := range filters.Filters {
    if !utils.Contains(searchPermissions, permission) {
      errorMessage := fmt.Sprintf("you do not have permission to search field '%v'", permission)
      return nil, errors.New(errorMessage)
    }
  }

  enterprises, err := sm.es.GetEnterprisesByFilters(filters)
  if err != nil {
    return nil, err
  }

  getPermissions, err := pm.FilterEnterpriseGetPermissions(permissions)
  if err != nil {
    sm.logger.Debug("Error getting get permissions")
    return nil, err
  }

  var rEnterprises []map[string]interface{}
  for _, enterprise := range *enterprises {
    enterpriseMap, err := utils.StructToMap(enterprise)
    if err != nil {
      return nil, err
    }

    for key := range enterpriseMap {
      if !utils.Contains(getPermissions, key) {
        delete(enterpriseMap, key)
      }
    }

    rEnterprises = append(rEnterprises, enterpriseMap)
  }

  for filter := range filters.Filters {
    err = sm.um.UsePermissions(u, []string{"ent_search_"+filter})
    if err != nil {
      return nil, err
    }
  }

  return &rEnterprises, nil
}

func (sm *SearchManager) GetSuggestions(params *e.GetSuggestions, u *u.User) ([]string, error) {
  var usedPermissions []string
  for filter := range params.Filters {
    if !pm.HavePermission(u, fmt.Sprintf("ent_search_%v", filter)) {
      errorMessage := fmt.Sprintf("you do not have permission to search '%v' property", filter)
      return nil, errors.New(errorMessage)
    }
    permissionSearch := fmt.Sprintf("ent_search_%v", filter)
    permissionGet := fmt.Sprintf("ent_get_%v", filter)
    usedPermissions = append(usedPermissions, permissionSearch, permissionGet)
  }
  if !pm.HavePermission(u, fmt.Sprintf("ent_search_%v", params.Field)){
    errorMessage := fmt.Sprintf("you do not have permission to search '%v' property", params.Field)
    return nil, errors.New(errorMessage)
  }
  usedPermissions = append(usedPermissions, params.Field)

  allowedFields := map[string]bool{
    "name":        true,
    "fantasyName": true,
    "state":       true,
    "city":        true,
    "district":    true,
  }

  if !allowedFields[params.Field] {
    errorMessage := fmt.Sprintf("field '%v' not permitted", params.Field)
    return nil, errors.New(errorMessage)
  }

  suggestions, err := sm.es.GetSuggestions(params)
  if err != nil {
    return nil, err
  }

  err = sm.um.UsePermissions(u, usedPermissions)
  if err != nil {
    return nil, err
  }

  return suggestions, nil
}