package middleware

import (
	"encoding/json"
	"net/http"

	"goproject-bank/helpers"
)

// AdminMiddleware ensures JWT is valid and user is admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if !helpers.ValidateAdminToken(tokenString) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Admin authorization required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// UserMiddleware ensures JWT is valid (admin OR user)
func UserMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        userID := r.Header.Get("User-ID")  

        if userID == "" || !helpers.ValidateToken(userID, tokenString) {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{"message": "User authorization required"})
            return
        }
        next.ServeHTTP(w, r)
    })
}

