package middleware

import (
	"context" // <-- NEW: Added context package for setting values
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

func GenerateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := VerifyJWT(tokenString) // VerifyJWT returns the token or an error
		if err != nil {
			// CRITICAL FIX: Use StatusUnauthorized for token verification failures
			http.Error(w, "invalid token", http.StatusUnauthorized) 
			return
		}
		
		// 1. Check if token is valid and extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		
		// 2. Extract the user_id from claims (claims are often float64)
		userIDfloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "user id not found in token claims", http.StatusUnauthorized)
			return
		}
		
		// 3. Convert to uint
		userID := uint(userIDfloat)
		
		// 4. Create a new context with the userID value
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		
		// 5. Serve the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}