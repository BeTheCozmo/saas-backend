package user_routes

import (
	"fmt"
	"net/http"
	"strings"
	u "uller/src/user"

	"github.com/gin-gonic/gin"
)

func (ur *UserRoutes) createUserHandler(c *gin.Context) {
  var user u.User
  err := c.ShouldBindBodyWithJSON(&user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
  }
  
  token, err := ur.um.CreateUser(&user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
    return
  }

  c.Header("Authorization", fmt.Sprintf("Bearer %v", token))
  c.Header("Content-Type", "application/json")

  c.AbortWithStatusJSON(http.StatusCreated, map[string]string{
    "token": token,
  })
}

func (ur *UserRoutes) authenticateHandler(c *gin.Context) {
  var user u.User
  err := c.ShouldBindBodyWithJSON(&user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
  }
  var credential string
  if user.Email != "" {
    credential = user.Email
  } else {
    credential = user.Phone
  }
  token, err := ur.um.GetAuthorization(credential, user.Password)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
  }

  c.Header("Authorization", fmt.Sprintf("Bearer %v", token))
  c.Header("Content-Type", "application/json")

  c.AbortWithStatusJSON(http.StatusOK, map[string]string{
    "token": token,
  })
}

func (ur *UserRoutes) getUserDetails(c *gin.Context) {
  if c.Request.Header["Authorization"] == nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "missing authorization header"})
    return
  }
  token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
  credential := c.Param("credential")
  requester, err := ur.a.ValidateJWT(token)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
  }
  requester, err = ur.um.GetByCredential(requester.Email)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }

  user, err := ur.um.GetDetails(requester, credential)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusNotFound, err.Error())
    return
  }

  newToken, err := ur.a.GetAuthorizationToken(user)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
    return
  }

  c.Header("Authorization", "Bearer "+newToken)
  c.AbortWithStatusJSON(http.StatusOK, user)
}

func (ur *UserRoutes) validateUser(c *gin.Context) {
  if c.Request.Header["Authorization"] == nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "missing authorization header"})
    return
  }
  token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
  user, err := ur.a.ValidateJWT(token)
  if err != nil {
    c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"valid": false})
    return
  }
  if user != nil {
    c.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{"valid": true})
    return
  }
  c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{"valid": false})
}