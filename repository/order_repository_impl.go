package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
)

type OrderRepositoryImpl struct{}

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImpl{}
}

func (repository *OrderRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, orderId int) (entity.Order, error) {
	SQL := "select id, product_id, sold_stock, status from orders where id = ?"
	rows, err := tx.QueryContext(ctx, SQL, orderId)
	helper.PanicIfError(err)
	defer rows.Close()

	order := entity.Order{}
	if rows.Next() {
		err := rows.Scan(&order.Id, &order.ProductId, &order.SoldStock, &order.Status)
		helper.PanicIfError(err)

		return order, nil
	} else {
		return order, errors.New("order not found")
	}
}

func (repository *OrderRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, order entity.Order) entity.Order {
	SQL := "insert into orders(product_id, sold_stock, status) values (?, ?, ?)"
	result, err := tx.ExecContext(ctx, SQL, order.ProductId, order.SoldStock, order.Status)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	order.Id = int(id)
	return order
}

func (repository *OrderRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, order entity.Order) entity.Order {
	SQL := "update orders set sold_stock = ?, status = ? where id = ?"
	_, err := tx.ExecContext(ctx, SQL, order.SoldStock, order.Status, order.Id)
	helper.PanicIfError(err)

	return order
}
