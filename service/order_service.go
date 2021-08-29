package service

import (
	"context"
	"github.com/anhsbolic/go-order-service-race-condition/model/web"
)

type OrderService interface {
	Create(ctx context.Context, request web.OrderCreateRequest)
}
