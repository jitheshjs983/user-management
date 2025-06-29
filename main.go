package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm/conf/controllers"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Reading from system env variables.")
	}

	// Read environment variables
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		user, pass, host, port, name,
	)

	// Connect to DB
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	handler := &controllers.Handler{DB: db}
	http.HandleFunc("/register", handler.RegisterUser)
	http.HandleFunc("/login", handler.LoginUser)
	http.HandleFunc("/change-password", handler.PasswordChange)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
