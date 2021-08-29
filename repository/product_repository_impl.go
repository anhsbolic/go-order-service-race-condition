package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
)

type ProductRepositoryImpl struct {
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}

func (repository *ProductRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []entity.Product {
	SQL := "select id, name from products"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		product := entity.Product{}
		err := rows.Scan(&product.Id, &product.Name)
		helper.PanicIfError(err)
		products = append(products, product)
	}

	return products
}

func (repository *ProductRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, productId int) (entity.Product, error) {
	SQL := "select id, name from products where id = ?"
	rows, err := tx.QueryContext(ctx, SQL, productId)
	helper.PanicIfError(err)
	defer rows.Close()

	product := entity.Product{}
	if rows.Next() {
		err := rows.Scan(&product.Id, &product.Name)
		helper.PanicIfError(err)

		return product, nil
	} else {
		return product, errors.New("product not found")
	}
}

func (repository *ProductRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product entity.Product) entity.Product {
	SQL := "insert into products(name) values (?)"
	result, err := tx.ExecContext(ctx, SQL, product.Name)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	product.Id = int(id)
	return product
}
