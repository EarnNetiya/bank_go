package middleware

import (
	"encoding/json"
	"net/http"

	"goproject-bank/helpers"
)

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