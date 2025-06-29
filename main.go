package main

import (
	"log"
	"net/http"

	"gorm/conf/controllers"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4"
var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

func main() {
	handler := &controllers.Handler{DB: db}
	http.HandleFunc("/user", handler.RegisterUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
