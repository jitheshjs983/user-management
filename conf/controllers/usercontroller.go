package controllers

import (
	"encoding/json"
	"fmt"
	"gorm/conf/middleware"
	"gorm/conf/models"
	"gorm/conf/service"
	"gorm/conf/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
		return
	}

	var user models.Users

	// Decode JSON body into user struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	emailExists, err := models.GetUserByEmail(h.DB, user.Email)
	if err != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}
	if emailExists {
		http.Error(w, "Email is already registered. Please try again with a different email.,", http.StatusConflict)
		return
	}
	phoneExists, err := models.GetUserByMobile(h.DB, user.Username)
	if err != nil {
		http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
		return
	}
	if phoneExists {
		http.Error(w, "mobile number is already registered. Please try again with a different mobile number.,", http.StatusConflict)
		return
	}
	if user.Password != "" {
		user.Password = models.HashPassword(user.Password)
	}

	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Println("DB error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	fmt.Println("✅ User added:", user)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
		return
	}

	var input models.LoginInput
	// Decode JSON body into user struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	loginType, err := utils.DetectLoginType(input.Login)
	if err != nil {
		http.Error(w, "Login must be a valid email or mobile number", http.StatusBadRequest)
		return
	}
	user, err := models.GetUserByLogin(h.DB, input.Login, loginType)
	if err != nil {
		log.Println("User lookup error:", err)
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}
	log.Printf("User found: %+v\n", user)
	// Compare password with bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.CreateToken(user.Username, user.Email, user.FirstName, user.LastName)
	if err != nil {
		panic(err)
	}
	resp := map[string]string{
		"message":  "Login successful",
		"username": user.Username,
		"email":    user.Email,
		"token":    tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) PasswordChange(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only POST method is allowed"))
		return
	}

	var input models.LoginInput
	// Decode JSON body into user struct
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	loginType, err := utils.DetectLoginType(input.Login)
	if err != nil {
		http.Error(w, "Login must be a valid email or mobile number", http.StatusBadRequest)
		return
	}
	user, err := models.GetUserByLogin(h.DB, input.Login, loginType)
	if err != nil {
		log.Println("User lookup error:", err)
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}
	// Compare password with bcrypt hash
	if models.IsSamePassword(user.Password, input.Password) {
		http.Error(w, "Please use a different password", http.StatusBadRequest)
		return
	}
	hashed := models.HashPassword(input.Password)
	if err := h.DB.Model(&user).Update("password", hashed).Error; err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{
		"message":  "Password updation successful",
		"username": user.Username,
		"email":    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) DashboardData(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"message": "Authentication Successful",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
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

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		http.Error(w, "Token missing exp", http.StatusUnauthorized)
		return
	}
	expiry := time.Unix(int64(expFloat), 0)

	middleware.BlacklistToken(jti, expiry)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}

func (h *Handler) GetNameFromPan(w http.ResponseWriter, r *http.Request) {
	var pan models.PanInput
	if err := json.NewDecoder(r.Body).Decode(&pan); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Call the service and get the result
	result, err := service.GetNameFromPan(pan.Pan)
	if err != nil {
		http.Error(w, "Failed to fetch PAN details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
