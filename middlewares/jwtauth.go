package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/swarnimcodes/apigo/utils"
)

var secretKey = os.Getenv("BEARER_TOKEN")

func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			message := "No authorization header sent"
			statusCode := http.StatusUnauthorized
			utils.SendErrorResponse(w, message, statusCode)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			message := "Invalid Authorization header format"
			statusCode := http.StatusUnauthorized
			utils.SendErrorResponse(w, message, statusCode)
			return
		}

		// tokenString := parts[1]

		// TODO: fix secret key
		secretKey := os.Getenv("BEARER_TOKEN")

		token, err := jwt.Parse(secretKey, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			fmt.Println(token)
			message := "Invalid JWT sent"
			statusCode := http.StatusUnauthorized
			utils.SendErrorResponse(w, message, statusCode)
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
