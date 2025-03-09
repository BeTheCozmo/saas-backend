package utils

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

func Contains(arr []string, value string) bool {
  for _, v := range arr {
    if v == value {
      return true
    }
  }
  return false
}

func HashPassword(password string) (string, error) {
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
  return string(hashedPassword), err
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
  jsonData, err := json.Marshal(obj)
    if err != nil {
        return nil, err
    }

    var result map[string]interface{}
    err = json.Unmarshal(jsonData, &result)
    if err != nil {
        return nil, err
    }

    return result, nil
}