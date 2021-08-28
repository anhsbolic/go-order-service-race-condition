package app

import (
	"database/sql"
	"github.com/anhsbolic/go-order-service-race-condition/controller"
	"github.com/anhsbolic/go-order-service-race-condition/exception"
	"github.com/anhsbolic/go-order-service-race-condition/repository"
	"github.com/anhsbolic/go-order-service-race-condition/service"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(db *sql.DB, validate *validator.Validate) *httprouter.Router {
	router := httprouter.New()

	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository, db, validate)
	productController := controller.NewProductController(productService)

	router.GET("/api/products", productController.FindAll)
	router.POST("/api/products", productController.Create)
	router.GET("/api/products/:productId", productController.FindById)

	router.PanicHandler = exception.ErrorHandler

	return router
}
