package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Uid      int64
	Nickname string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Phone    string `json:"phone"`
	jwt.RegisteredClaims
}

// CreateToken generate token
func CreateToken(c *Claims, key string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedString, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("generate token failed" + err.Error())
	}
	return signedString, nil
}
