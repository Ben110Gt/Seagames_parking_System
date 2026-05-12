package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("SECRET")

func GenerateToken(userID string, UserName string, role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": UserName,
		"role":      role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtSecret)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
}
