package service

import (
	"context"
	"github.com/anhsbolic/go-order-service-race-condition/model/web"
)

type ProductService interface {
	FindAll(ctx context.Context) []web.ProductResponse
	FindById(ctx context.Context, productId int) web.ProductResponse
	Create(ctx context.Context, request web.ProductCreateRequest) web.ProductResponse
}
