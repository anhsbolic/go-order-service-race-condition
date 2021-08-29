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

	// init repositories
	inventoryRepository := repository.NewInventoryRepository()
	orderRepository := repository.NewOrderRepository()
	productRepository := repository.NewProductRepository()

	// init product service & controller
	productService := service.NewProductService(inventoryRepository, productRepository, db, validate)
	productController := controller.NewProductController(productService)

	// init order service & controller
	orderService := service.NewOrderService(inventoryRepository, orderRepository, productRepository, db, validate)
	orderController := controller.NewOrderController(orderService)

	// routes : products
	router.GET("/api/products", productController.FindAll)
	router.POST("/api/products", productController.Create)
	router.GET("/api/products/:productId", productController.FindById)
	router.GET("/api/products/:productId/inventory", productController.FindInventoryByProductId)
	router.PUT("/api/products/:productId/inventory", productController.UpdateInventoryByProductId)

	// routes : orders
	router.POST("/api/orders", orderController.Create)

	// router handler
	router.PanicHandler = exception.ErrorHandler

	return router
}
