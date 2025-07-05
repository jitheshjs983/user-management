package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims struct (customize as needed)
type MyClaims struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

// CreateToken generates a JWT token string
func CreateToken(username, email, firstName, lastName string) (string, error) {
	now := time.Now()
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	claims := MyClaims{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "User Management",
		},
	}
	fmt.Println("JWT_SECRET:", jwtSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
