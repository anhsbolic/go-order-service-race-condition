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

type ProductServiceImpl struct {
	InventoryRepository repository.InventoryRepository
	ProductRepository   repository.ProductRepository
	DB                  *sql.DB
	Validate            *validator.Validate
}

func NewProductService(
	inventoryRepository repository.InventoryRepository,
	productRepository repository.ProductRepository,
	db *sql.DB, validate *validator.Validate,
) ProductService {
	return &ProductServiceImpl{
		InventoryRepository: inventoryRepository,
		ProductRepository:   productRepository,
		DB:                  db,
		Validate:            validate,
	}
}

func (service *ProductServiceImpl) FindAll(ctx context.Context) []web.ProductResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	products := service.ProductRepository.FindAll(ctx, tx)

	return web.ToProductResponses(products)
}

func (service *ProductServiceImpl) FindById(ctx context.Context, productId int) web.ProductResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindById(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return web.ToProductResponse(product)
}

func (service *ProductServiceImpl) Create(ctx context.Context, request web.ProductCreateRequest) web.ProductResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// create new product
	product := entity.Product{
		Name: request.Name,
	}
	product = service.ProductRepository.Save(ctx, tx, product)

	// create inventory product
	inventory := entity.Inventory{
		ProductId:      product.Id,
		StoredStock:    0,
		AvailableStock: 0,
		ReservedStock:  0,
	}
	inventory = service.InventoryRepository.Save(ctx, tx, inventory)

	return web.ToProductResponse(product)
}

func (service *ProductServiceImpl) FindInventoryByProductId(ctx context.Context, productId int) web.InventoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	product, err := service.ProductRepository.FindById(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	inventory, err := service.InventoryRepository.FindByProductId(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return web.InventoryResponse{
		ProductName:    product.Name,
		StoredStock:    inventory.StoredStock,
		AvailableStock: inventory.AvailableStock,
		ReservedStock:  inventory.ReservedStock,
	}
}

func (service *ProductServiceImpl) UpdateInventoryByProductId(ctx context.Context, productId int, request web.InventoryUpdateRequest) web.InventoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// find product
	product, err := service.ProductRepository.FindById(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// find inventory product
	inventory, err := service.InventoryRepository.FindByProductId(ctx, tx, productId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// update stored stock & available stock
	inventory.StoredStock = inventory.StoredStock + request.Stock
	inventory.AvailableStock = inventory.AvailableStock + request.Stock

	// update inventory
	inventory = service.InventoryRepository.Update(ctx, tx, inventory)

	return web.InventoryResponse{
		ProductName:    product.Name,
		StoredStock:    inventory.StoredStock,
		AvailableStock: inventory.AvailableStock,
		ReservedStock:  inventory.ReservedStock,
	}
}
