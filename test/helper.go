package test

import (
	"database/sql"
	"fmt"
	"github.com/anhsbolic/go-order-service-race-condition/app"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func SetupTestDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		helper.PanicIfError(err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	db, err := sql.Open("mysql", source)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func SetupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	router := app.NewRouter(db, validate)

	return router
}

func ClearProductsTable(db *sql.DB) {
	_, err := db.Exec("delete from products")
	if err != nil {
		fmt.Println(err)
	}
}

func ClearOrdersTable(db *sql.DB) {
	_, err := db.Exec("delete from orders")
	if err != nil {
		fmt.Println(err)
	}
}