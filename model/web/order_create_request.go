package web

type OrderCreateRequest struct {
	ProductId int `validate:"required" json:"product_id"`
	Total     int `validate:"required,min=0" json:"total"`
}
