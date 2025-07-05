package main

import (
	"log"
	"net/http"

	"gorm/conf/controllers"
	"gorm/conf/database"
	"gorm/conf/middleware"
	"gorm/conf/utils"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	utils.InitRedis()
	database.Init()

	handler := &controllers.Handler{DB: database.DB}

	http.HandleFunc("/register", handler.RegisterUser)
	http.HandleFunc("/login", handler.LoginUser)
	http.HandleFunc("/change-password", handler.PasswordChange)
	http.Handle("/dashboard", middleware.JWTMiddleware(http.HandlerFunc(handler.DashboardData)))
	http.Handle("/logout", middleware.JWTMiddleware(http.HandlerFunc(handler.LogoutUser)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
