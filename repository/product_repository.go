package repository

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
)

type ProductRepository interface {
	FindAll(ctx context.Context, tx *sql.Tx) []entity.Product
	FindById(ctx context.Context, tx *sql.Tx, productId int) (entity.Product, error)
	Save(ctx context.Context, tx *sql.Tx, product entity.Product) entity.Product
}
