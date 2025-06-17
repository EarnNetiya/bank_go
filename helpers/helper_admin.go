package helpers

import (
	"encoding/json"
	"goproject-bank/interfaces"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidationAdmin(values []interfaces.Validation) bool {
	username := regexp.MustCompile(`^([A-Za-z0-9]{5,})+$`)
	email := regexp.MustCompile(`^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z]+$`)

	for i := 0; i < len(values); i++ {
		switch values[i].Valid {
		case "username":
			if !username.MatchString(values[i].Value) {
				return false
			}
		case "email":
			if !email.MatchString(values[i].Value) {
				return false
			}
		case "password":
			if len(values[i].Value) < 5 {
				return false
			}
		}
	}
	return true
}

func PanicHandlerAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				resp := interfaces.ErrResponse{Message: "Internal server error"}
				json.NewEncoder(w).Encode(resp)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ValidateAdminToken(tokenString string) bool {

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        return JwtSecret, nil // ใช้ secret key เดียวกัน
    })

    if err != nil || !token.Valid {
        return false
    }

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid token claims")
		return false
	}

	if isAdmin, ok := claims["admin"].(bool); !ok || !isAdmin {
		log.Println("Admin flag missing or false")
		return false
	}

	if code, ok := claims["verification_code"].(string); !ok || code != AdminVerificationCode {
		log.Println("Verification code invalid")
		return false
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) <= time.Now().Unix() {
		log.Println("Token expired")
		return false
	}

	return true
}



func ValidateToken(userID, tokenString string) bool {
	return true
}
