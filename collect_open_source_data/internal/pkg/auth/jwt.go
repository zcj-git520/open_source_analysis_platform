package auth

import (
	"context"
	"errors"
	"fmt"
	jwtauth "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Uid      int64  `json:"uid"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Phone    string `json:"phone"`
	jwt.RegisteredClaims
}

// CreateToken generate token
func CreateToken(c Claims, key string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedString, err := claims.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("generate token failed" + err.Error())
	}
	return signedString, nil
}

func ParsToken(tokenString, key string) *Claims {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		fmt.Println("解析 token 失败：", err)
		return nil
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims
	}
	return nil

}

func UserInfo(ctx context.Context) *Claims {
	return &Claims{
		Uid:      GetUid(ctx),
		Nickname: GetNickname(ctx),
		Email:    GetEmail(ctx),
		Avatar:   GetAvatar(ctx),
		Gender:   GetGender(ctx),
		Phone:    GetPhone(ctx),
	}
}

func GetUid(ctx context.Context) int64 {
	var uId int64
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		i, ok := c["uid"].(float64)
		if !ok {
			return uId
		}
		uId = int64(i)
	}
	return uId
}

func GetNickname(ctx context.Context) string {
	var nickname string
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		nickname = c["username"].(string)
	}
	return nickname
}

func GetEmail(ctx context.Context) string {
	var email string
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		email = c["email"].(string)
	}
	return email
}

func GetAvatar(ctx context.Context) string {
	var avatar string
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		avatar = c["avatar"].(string)
	}
	return avatar
}

func GetGender(ctx context.Context) int {
	var gender int
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		i, ok := c["gender"].(float64)
		if !ok {
			return gender
		}
		gender = int(i)
	}
	return gender
}

func GetPhone(ctx context.Context) string {
	var phone string
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		phone = c["phone"].(string)
	}
	return phone
}
