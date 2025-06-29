package controllers

import (
	"encoding/json"
	"fmt"
	"gorm/conf/models"
	"log"
	"net/http"

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
