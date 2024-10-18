package auth

import (
	"context"
	"errors"
	"fmt"
	jwtauth "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
	"strings"
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

func passTokenString(ctx context.Context) *Claims {
	authorizationKey := "Authorization"
	bearerWord := "Bearer"
	if header, ok := transport.FromServerContext(ctx); ok {
		auths := strings.SplitN(header.RequestHeader().Get(authorizationKey), " ", 2)
		if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
			return nil
		}
		return ParsToken(auths[1], "hqFr%3ddt32DGlSTOI5cO6@TH#fFwYnP$S")
	}
	return nil
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
	} else {
		if info := passTokenString(ctx); info != nil {
			uId = info.Uid
		}
	}
	return uId
}

func GetNickname(ctx context.Context) string {
	var nickname string
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		nickname = c["username"].(string)
	} else {
		if info := passTokenString(ctx); info != nil {
			nickname = info.Nickname
		}
	}
	return nickname
}

func GetEmail(ctx context.Context) string {
	var email string
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		email = c["email"].(string)
	} else {
		if info := passTokenString(ctx); info != nil {
			email = info.Email
		}
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
	} else {
		if info := passTokenString(ctx); info != nil {
			gender = info.Gender
		}
	}
	return gender
}

func GetPhone(ctx context.Context) string {
	var phone string
	if claims, ok := jwtauth.FromContext(ctx); ok {
		c := claims.(jwt.MapClaims)
		phone = c["phone"].(string)
	} else {
		if info := passTokenString(ctx); info != nil {
			phone = info.Phone
		}
	}
	return phone
}
