package main

import (
	"github.com/anhsbolic/go-order-service-race-condition/app"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	db := app.NewDB()
	validate := validator.New()
	router := app.NewRouter(db, validate)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
