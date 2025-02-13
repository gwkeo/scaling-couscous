package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gwkeo/scaling-couscous/internal/app/repo/models"
	"strings"
)

type Claims struct {
	id   int64
	role models.Role
}

func GenerateToken(secret string, id int64, role models.Role) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": role,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret string) (*jwt.Token, error) {
	t := strings.Split(tokenString, " ")
	if len(t) != 2 {
		return nil, errors.New("token format error")
	}
	tokenString = t[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
