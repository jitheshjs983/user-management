package middleware

import (
	"fmt"
	"gorm/conf/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var tokenBlacklist = make(map[string]time.Time)

// Add token to blacklist
func BlacklistToken(jti string, expiry time.Time) error {
	duration := time.Until(expiry)
	return utils.RedisClient.Set(utils.Ctx, jti, "revoked", duration).Err()
}

// Check if token is blacklisted
func isBlacklisted(jti string) bool {
	val, err := utils.RedisClient.Get(utils.Ctx, jti).Result()
	if err != nil {
		return false
	}
	return val == "revoked"
}
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok {
			http.Error(w, "Token missing jti", http.StatusUnauthorized)
			return
		}

		if isBlacklisted(jti) {
			http.Error(w, "Token has been revoked", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
