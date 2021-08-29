package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
)

type InventoryRepositoryImpl struct {
}

func NewInventoryRepository() InventoryRepository {
	return &InventoryRepositoryImpl{}
}

func (repository *InventoryRepositoryImpl) FindByProductId(ctx context.Context, tx *sql.Tx, productId int) (entity.Inventory, error) {
	SQL := "select product_id, stored_stock, available_stock, reserved_stock from inventories where product_id = ?"
	rows, err := tx.QueryContext(ctx, SQL, productId)
	helper.PanicIfError(err)
	defer rows.Close()

	inventory := entity.Inventory{}
	if rows.Next() {
		err := rows.Scan(&inventory.ProductId, &inventory.StoredStock, &inventory.AvailableStock, &inventory.ReservedStock)
		helper.PanicIfError(err)

		return inventory, nil
	} else {
		return inventory, errors.New("inventory not found")
	}
}

func (repository *InventoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, inventory entity.Inventory) entity.Inventory {
	SQL := "insert into inventories(product_id, stored_stock, available_stock, reserved_stock) values (?, ?, ?, ?)"
	_, err := tx.ExecContext(ctx, SQL, inventory.ProductId, inventory.StoredStock, inventory.AvailableStock, inventory.ReservedStock)
	helper.PanicIfError(err)

	return inventory
}

func (repository *InventoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, inventory entity.Inventory) entity.Inventory {
	SQL := "update inventories set available_stock = ?, reserved_stock = ? where product_id = ?"
	_, err := tx.ExecContext(ctx, SQL, inventory.AvailableStock, inventory.ReservedStock, inventory.ProductId)
	helper.PanicIfError(err)

	return inventory
}
