package web

import "github.com/anhsbolic/go-order-service-race-condition/model/entity"

type ProductResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func ToProductResponse(product entity.Product) ProductResponse {
	return ProductResponse{
		Id:   product.Id,
		Name: product.Name,
	}
}

func ToProductResponses(products []entity.Product) []ProductResponse {
	var productResponses []ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, ToProductResponse(product))
	}
	return productResponses
}
