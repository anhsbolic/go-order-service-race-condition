package repository

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
)

type OrderRepository interface {
	FindById(ctx context.Context, tx *sql.Tx, orderId int) (entity.Order, error)
	Save(ctx context.Context, tx *sql.Tx, order entity.Order) entity.Order
	Update(ctx context.Context, tx *sql.Tx, order entity.Order) entity.Order
}
