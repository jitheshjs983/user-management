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
	user := models.Users{
		FirstName: "Jithesh",
		LastName:  "Jose",
		Username:  "8590811971",
		Password:  models.HashPassword("Test@123"),
		Email:     "jitheshjs983@gmail.com",
	}
	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		log.Println("DB error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	fmt.Println("✅ Fields added:", user)
}
