package web

type InventoryResponse struct {
	ProductName    string `json:"product_name"`
	StoredStock    int    `json:"stored_stock"`
	AvailableStock int    `json:"available_stock"`
	ReservedStock  int    `json:"reserved_stock"`
}
