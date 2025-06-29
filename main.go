package main

import (
	"log"
	"net/http"

	"gorm/conf/controllers"
	"gorm/conf/database"
)

func main() {
	database.Init()

	handler := &controllers.Handler{DB: database.DB}

	http.HandleFunc("/register", handler.RegisterUser)
	http.HandleFunc("/login", handler.LoginUser)
	http.HandleFunc("/change-password", handler.PasswordChange)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
