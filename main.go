package main

import (
	"log"
	"net/http"

	"gorm/conf/controllers"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=true"
var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

func main() {
	handler := &controllers.Handler{DB: db}
	http.HandleFunc("/register", handler.RegisterUser)
	http.HandleFunc("/login", handler.LoginUser)
	http.HandleFunc("/change-password", handler.PasswordChange)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
