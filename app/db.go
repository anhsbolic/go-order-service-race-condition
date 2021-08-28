package app

import (
	"database/sql"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"time"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:newpass@tcp(localhost:3306)/order-service-race-condition-poc")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
