package helpers

import (
	// "crypto/md5"
	// "encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"goproject-bank/interfaces"
	"log"
	"net/http"
	"regexp"


	"github.com/golang-jwt/jwt/v5"
	// "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)

	return string(hashed)
}


func Validation(values []interfaces.Validation) bool {
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

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func ()  {
			error := recover()
			if error != nil {
				log.Println(error)
				resp := interfaces.ErrResponse{Message: "Internal server error"}
				json.NewEncoder(w).Encode(resp)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ValidateUserToken(userId string, tokenString string) bool {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })

    if err != nil || !token.Valid {
        log.Println("Token parse error:", err)
        return false
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return false
    }

    uid := fmt.Sprintf("%v", claims["user_id"])
    return uid == userId
}

var secretKey = []byte("JwtSecret")

func ExtractUserID(tokenString string) (uint, error) {
    if tokenString == "" {
        return 0, errors.New("token is empty")
    }

    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secretKey, nil
    })
    if err != nil || !token.Valid {
        return 0, fmt.Errorf("invalid token: %w", err)
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("invalid claims")
    }

    userIDFloat, ok := claims["user_id"].(float64)
    if !ok {
        return 0, errors.New("user_id not found in token")
    }

    return uint(userIDFloat), nil
}

func ExtractUserIDFromToken(tokenString string) (string, error) {
    if tokenString == "" {
        return "", fmt.Errorf("token is empty")
    }
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        return JwtSecret, nil 
    })

    if err != nil || !token.Valid {
        return "", fmt.Errorf("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", fmt.Errorf("invalid claims")
    }

    userID := fmt.Sprintf("%v", claims["user_id"])
    if userID == "" {
        return "", fmt.Errorf("user_id not found")
    }

    return userID, nil
}

func GetUserIDFromToken(tokenString string) uint {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			return uint(userIDFloat)
		}
	}

	log.Println("GetUserIDFromToken failed:", err)
	return 0
}

func ExtractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return authHeader
}

func ExtractTokenFromRequest(r *http.Request) string {
    authHeader := r.Header.Get("Authorization")
    log.Println("Authorization header:", authHeader) // Debug
    if authHeader == "" {
        return ""
    }
    if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
        token := authHeader[7:]
        log.Println("Extracted token:", token) // Debug
        return token
    }
    log.Println("No Bearer prefix, using raw header:", authHeader) // Debug
    return authHeader
}