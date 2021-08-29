package repository

import (
	"context"
	"database/sql"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
)

type InventoryRepository interface {
	FindByProductId(ctx context.Context, tx *sql.Tx, productId int) (entity.Inventory, error)
	Save(ctx context.Context, tx *sql.Tx, inventory entity.Inventory) entity.Inventory
	Update(ctx context.Context, tx *sql.Tx, inventory entity.Inventory) entity.Inventory
}
