package app

import (
	"database/sql"
	"fmt"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func NewDB() *sql.DB {
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
