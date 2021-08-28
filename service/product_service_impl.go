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
	ProductRepository repository.ProductRepository
	DB                *sql.DB
	Validate          *validator.Validate
}

func NewProductService(productRepository repository.ProductRepository, db *sql.DB, validate *validator.Validate) ProductService {
	return &ProductServiceImpl{
		ProductRepository: productRepository,
		DB:                db,
		Validate:          validate,
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

	product := entity.Product{
		Name: request.Name,
	}

	product = service.ProductRepository.Save(ctx, tx, product)

	return web.ToProductResponse(product)
}
