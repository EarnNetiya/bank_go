package helpers

import (
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func ValidateAdminToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		adminVal, ok := claims["admin"].(bool)
		if ok && adminVal {
			exp, ok := claims["exp"].(float64)
			if !ok || int64(exp) < time.Now().Unix() {
				return false
			}
			return true
		}
	}
	return false
}


