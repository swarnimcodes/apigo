package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("Bearer token123")

func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization Header Missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		expTime := token.Claims.(jwt.MapClaims)["exp"].(float64)
		expDateTime := time.Unix(int64(expTime), 0)
		remainingDuration := expDateTime.Sub(time.Now())
		log.Printf("Token will expire at: %s\n", expDateTime)
		log.Printf("Remaining time until token expiration (seconds): %f\n", remainingDuration.Seconds())
		next.ServeHTTP(w, r)
	})
}
