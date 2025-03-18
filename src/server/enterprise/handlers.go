package enterprise_routes

import (
	"net/http"
	"strings"
	e "uller/src/search/enterprise"

	"github.com/gin-gonic/gin"
)

func (er *EnterpriseRoutes) getEnterprises(c *gin.Context) {
  if c.Request.Header["Authorization"] == nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "missing authorization header"})
    return
  }
  token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
  user, err := er.a.ValidateJWT(token)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
    return
  }
  user, err = er.um.GetByCredential(user.Email)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }

  var rawRequest map[string]interface{}
  if err := c.ShouldBindJSON(&rawRequest); err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
    return
  }

  ordersRaw, _ := rawRequest["orders"].(map[string]interface{})
  if ordersRaw == nil {
    ordersRaw = make(map[string]interface{})
  }
  filters, _ := rawRequest["filters"].(map[string]interface{})
  if filters != nil {
    for k, v := range filters {
      er.logger.Debug("Filter:", k)
      switch v := v.(type) {
      case []interface{}:
        arr := make([]string, 0, len(v))
        for _, item := range v {
          if str, ok := item.(string); ok {
            arr = append(arr, str)
          } else {
            er.logger.Warn("Non-string item in array filter", k, item)
          }
        }
        er.logger.Debug("Is array", arr)
        filters[k] = arr
      case map[string]interface{}:
        numMap := make(map[string]float64)
        for key, val := range v {
          if num, ok := val.(float64); ok {
            numMap[key] = num
          } else {
            er.logger.Warn("Non-numeric value in object filter", k, key, val)
          }
        }
        er.logger.Debug("Is map:", numMap)
        filters[k] = numMap
      default:
        er.logger.Warn("Unsupported filter type for", k, v)
      }
    }
  } else {
    filters = make(map[string]interface{})
  }

  er.logger.Debug(filters)

  pageRaw, ok := rawRequest["page"].(uint64)
  if !ok {
    pageRaw = 1
  }

  orders := make(map[string]string)
  for k, v := range ordersRaw {
    if strVal, ok := v.(string); ok && (strVal == "asc" || strVal == "desc") {
      orders[k] = strVal
    }
  }

  params := &e.GetEnterprisesFilter{
    Orders:  orders,
    Filters: filters,
    Page:    pageRaw,
  }

  er.logger.Debug("Getting enterprises", params)
  enterprises, err := er.sm.GetEnterprisesByFilters(params, user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
    return
  }

  newToken, err := er.a.GetAuthorizationToken(user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }
  c.Header("Authorization", "Bearer "+newToken)
  c.AbortWithStatusJSON(http.StatusOK, enterprises)
}

func (er *EnterpriseRoutes) getByDocument(c *gin.Context) {
  if c.Request.Header["Authorization"] == nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "missing authorization header"})
    return
  }
  token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
  document := c.Param("document")
  user, err := er.a.ValidateJWT(token)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
    return
  }
  user, err = er.um.GetByCredential(user.Email)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }
  enterprise, err := er.sm.GetEnterpriseBaseByDocumentMap(document, user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
    return
  }

  newToken, err := er.a.GetAuthorizationToken(user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }
  c.Header("Authorization", "Bearer "+newToken)
  c.AbortWithStatusJSON(http.StatusOK, enterprise)
}

func (er *EnterpriseRoutes) getHistory(c *gin.Context) {
  if c.Request.Header["Authorization"] == nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "missing authorization header"})
    return
  }
  token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
  document := c.Param("document")
  user, err := er.a.ValidateJWT(token)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
    return
  }
  user, err = er.um.GetByCredential(user.Email)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }
  enterprises, err := er.sm.GetEnterpriseHistoryMap(document, 20, user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
    return
  }

  newToken, err := er.a.GetAuthorizationToken(user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }
  c.Header("Authorization", "Bearer "+newToken)
  c.AbortWithStatusJSON(http.StatusOK, enterprises)
}

func (er *EnterpriseRoutes) getFilterSuggestions(c *gin.Context) {
  if c.Request.Header["Authorization"] == nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "missing authorization header"})
    return
  }
  token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
  user, err := er.a.ValidateJWT(token)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
    return
  }

  er.logger.Info("Received request for filter suggestions")

  var params e.GetSuggestions
  if err := c.ShouldBindJSON(&params); err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
    return
  }

  suggestions, err := er.sm.GetSuggestions(&params, user)
  if err != nil {
    er.logger.Debug(err)
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
    return
  }
  er.logger.Debug("Request params:", params)

  newToken, err := er.a.GetAuthorizationToken(user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }
  c.Header("Authorization", "Bearer "+newToken)
  c.AbortWithStatusJSON(http.StatusOK, suggestions)
}