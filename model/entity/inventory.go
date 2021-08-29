package entity

type Inventory struct {
	ProductId      int `json:"product_id"`
	StoredStock    int `json:"stored_stock"`
	AvailableStock int `json:"available_stock"`
	ReservedStock  int `json:"reserved_stock"`
}
