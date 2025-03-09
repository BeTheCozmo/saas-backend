package authorization

import (
	"errors"
	"time"
	p "uller/src/permission"
	pl "uller/src/plan"
	r "uller/src/role"
	"uller/src/user"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Authorization struct {
  jwtSecret string
}

func New() *Authorization {
  return &Authorization{}
}

func (a *Authorization) Configure(jwtSecret string) {
  a.jwtSecret = jwtSecret
}

func (a *Authorization) serializeUser(user *user.User) (string, error) {

  expirationTime := time.Now().Add(time.Hour * 24 * 7).Unix()

  claims := jwt.MapClaims{
    "name":        user.Name,
    "email":       user.Email,
    "phone":       user.Phone,
    "permissions": user.Permissions,
    "role":        user.Role,
    "plan":        user.Plan,
    "exp":         expirationTime,
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  if token == nil {
    return "", errors.New("failed to create JWT token")
  }
  return token.SignedString([]byte(a.jwtSecret))
}

func (a *Authorization) GetAuthorizationToken(user *user.User) (string, error) {
  return a.serializeUser(user)
}

func (a *Authorization) ValidateJWT(tokenString string) (*user.User, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("assinatura inválida")
    }
    return []byte(a.jwtSecret), nil
  })

  if err != nil {
    return nil, err
  }

  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    if exp, ok := claims["exp"].(float64); ok {
      expTime := time.Unix(int64(exp), 0)
      if time.Now().After(expTime) {
        return nil, errors.New("expired token")
      }
    }

    var user user.User
    err := a.mapClaimsToStruct(claims, &user)
    return &user, err
  }

  return nil, errors.New("invalid token")
}

func (a *Authorization) mapClaimsToStruct(claims jwt.MapClaims, target *user.User) error {
  // A biblioteca jwt fornece uma forma de mapear as claims, mas precisamos garantir
  // que estamos tratando todos os tipos corretos.
  if name, ok := claims["name"].(string); ok {
    target.Name = name
  } else {
    return errors.New("invalid field 'name'")
  }

  if email, ok := claims["email"].(string); ok {
    target.Email = email
  } else {
    return errors.New("invalid field 'email'")
  }

  if phone, ok := claims["phone"].(string); ok {
    target.Phone = phone
  } else {
    return errors.New("invalid field 'phone'")
  }

  if permissions, ok := claims["permissions"].([]interface{}); ok {
    for _, perm := range permissions {
      if permMap, ok := perm.(map[string]interface{}); ok {
        var permission p.Permission
        if name, ok := permMap["name"].(string); ok {
          permission.Name = name
        }
        if isLimited, ok := permMap["isLimited"].(bool); ok {
          permission.IsLimited = isLimited
        }
        if remainingUses, ok := permMap["remainingUses"].(float64); ok {
          permission.RemainingUses = uint32(remainingUses)
        }
        target.Permissions = append(target.Permissions, &permission)
      } else {
        return errors.New("campo 'permissions' contém dados em formato inválido")
      }
    }
  } else {
    return errors.New("campo 'permissions' não encontrado ou não é um array")
  }

  if role, ok := claims["role"].(map[string]interface{}); ok {
    var rRole r.Role

    // Mapeia os campos de 'role'
    if name, ok := role["name"].(string); ok {
      rRole.Name = name
    } else {
      return errors.New("campo 'role.name' inválido")
    }

    // Verifica se existe o campo de permissões dentro de role
    if permissions, ok := role["permissions"].([]interface{}); ok {
      for _, perm := range permissions {
        if permMap, ok := perm.(map[string]interface{}); ok {
          var permission p.Permission
          if name, ok := permMap["name"].(string); ok {
            permission.Name = name
          }
          if isLimited, ok := permMap["isLimited"].(bool); ok {
            permission.IsLimited = isLimited
          }
          if remainingUses, ok := permMap["remainingUses"].(float64); ok {
            permission.RemainingUses = uint32(remainingUses)
          }
          rRole.Permissions = append(rRole.Permissions, &permission)
        } else {
          return errors.New("campo 'role.permissions' contém dados em formato inválido")
        }
      }
    }

    target.Role = &rRole
  } else {
    return errors.New("campo 'role' inválido ou não encontrado")
  }

  if plan, ok := claims["plan"].(map[string]interface{}); ok {
    var pPlan pl.Plan

    // Mapeia os campos de 'plan'
    if name, ok := plan["name"].(string); ok {
      pPlan.Name = name
    } else {
      return errors.New("campo 'plan.name' inválido")
    }

    if permissions, ok := plan["permissions"].([]interface{}); ok {
      for _, perm := range permissions {
        if permMap, ok := perm.(map[string]interface{}); ok {
          var permission p.Permission
          if name, ok := permMap["name"].(string); ok {
            permission.Name = name
          }
          if isLimited, ok := permMap["isLimited"].(bool); ok {
            permission.IsLimited = isLimited
          }
          if remainingUses, ok := permMap["remainingUses"].(float64); ok {
            permission.RemainingUses = uint32(remainingUses)
          }
          pPlan.Permissions = append(pPlan.Permissions, &permission)
        } else {
          return errors.New("campo 'plan.permissions' contém dados em formato inválido")
        }
      }
    }

    if startedAt, ok := plan["startedAt"].(float64); ok {
      pPlan.StartedAt = uint64(startedAt)
    } else {
      return errors.New("campo 'plan.startedAt' inválido")
    }

    if endAtBase, ok := plan["endAtBase"].(float64); ok {
      pPlan.EndAtBase = uint64(endAtBase)
    } else {
      return errors.New("campo 'plan.endAtBase' inválido")
    }

    if priceBase, ok := plan["priceBase"].(float64); ok {
      pPlan.PriceBase = uint64(priceBase)
    } else {
      return errors.New("campo 'plan.priceBase' inválido")
    }

    target.Plan = &pPlan
  } else {
    return errors.New("campo 'plan' inválido ou não encontrado")
  }

  return nil
}
