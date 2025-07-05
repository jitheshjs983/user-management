package main

import (
	"log"
	"net/http"

	"gorm/conf/controllers"
	"gorm/conf/database"
	"gorm/conf/middleware"
)

func main() {
	database.Init()

	handler := &controllers.Handler{DB: database.DB}

	http.HandleFunc("/register", handler.RegisterUser)
	http.HandleFunc("/login", handler.LoginUser)
	http.HandleFunc("/change-password", handler.PasswordChange)
	http.Handle("/dashboard", middleware.JWTMiddleware(http.HandlerFunc(handler.DashboardData)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
