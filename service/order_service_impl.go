package service

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/go-order-service-race-condition/exception"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
	"github.com/anhsbolic/go-order-service-race-condition/model/web"
	"github.com/anhsbolic/go-order-service-race-condition/repository"
	"github.com/go-playground/validator/v10"
)

type OrderServiceImpl struct {
	InventoryRepository repository.InventoryRepository
	OrderRepository     repository.OrderRepository
	ProductRepository   repository.ProductRepository
	DB                  *sql.DB
	Validate            *validator.Validate
}

func NewOrderService(
	inventoryRepository repository.InventoryRepository,
	orderRepository repository.OrderRepository,
	productRepository repository.ProductRepository,
	db *sql.DB,
	validate *validator.Validate,
) OrderService {
	return &OrderServiceImpl{
		InventoryRepository: inventoryRepository,
		OrderRepository:     orderRepository,
		ProductRepository:   productRepository,
		DB:                  db,
		Validate:            validate,
	}
}

func (service *OrderServiceImpl) Create(ctx context.Context, request web.OrderCreateRequest) {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// find product
	product, err := service.ProductRepository.FindById(ctx, tx, request.ProductId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// find inventory product
	inventory, err := service.InventoryRepository.FindByProductId(ctx, tx, product.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// available stock inventory validation
	if inventory.AvailableStock <= 0 {
		panic(exception.NewBadRequestError("Product Available Stock Is Empty"))
	}

	// create order
	order := entity.Order{
		ProductId: product.Id,
		SoldStock: request.Total,
		Status:    entity.OrderStatusCreated,
	}
	order = service.OrderRepository.Save(ctx, tx, order)

	// update inventory available & reserved stock
	inventory.AvailableStock = inventory.AvailableStock - request.Total
	inventory.ReservedStock = inventory.ReservedStock + request.Total
	inventory = service.InventoryRepository.Update(ctx, tx, inventory)
}
