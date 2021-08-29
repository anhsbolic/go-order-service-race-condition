package entity

var (
	OrderStatusCreated   = 1
	OrderStatusPaid      = 2
	OrderStatusCancelled = 3
)

type Order struct {
	Id        int `json:"id"`
	ProductId int `json:"product_id"`
	SoldStock int `json:"sold_stock"`
	Status    int `json:"status"`
}
