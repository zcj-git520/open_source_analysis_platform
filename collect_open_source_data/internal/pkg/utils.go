package pkg

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strconv"
)

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func ToString(data any) string {
	switch data.(type) {
	case string:
		return data.(string)
	case int:
		return strconv.Itoa(data.(int))
	case int32:
		return strconv.FormatInt(int64(data.(int32)), 10)
	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(data.(float64), 'f', -1, 32)
	case bool:
		if data.(bool) {
			return "true"
		}
		return "false"
	case []byte:
		return string(data.([]byte))
	default:
		buf, err := json.Marshal(data)
		if err != nil {
			return ""
		}
		return string(buf)
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
