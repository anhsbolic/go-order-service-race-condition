package web

type InventoryUpdateRequest struct {
	Stock int `validate:"required,min=0" json:"stock"`
}
