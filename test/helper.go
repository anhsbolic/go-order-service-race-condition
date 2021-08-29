package test

import (
	"database/sql"
	"fmt"
	"github.com/anhsbolic/go-order-service-race-condition/app"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

func SetupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:newpass@tcp(localhost:3306)/order-service-race-condition-poc")
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