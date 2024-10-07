package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip JWT Verification for OPTIONS Requests
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Split the header to get the token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// Retrieve the JWT secret from environment variables
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		if jwtSecret == nil {
			http.Error(w, "JWT secret is not set", http.StatusInternalServerError)
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Printf("Token parse error: %v", err)
			return
		}

		// Optionally, you can extract claims from the token and attach them to the request context
		next.ServeHTTP(w, r)
	})
}
