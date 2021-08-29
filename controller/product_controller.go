package controller

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ProductController interface {
	FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params)

	FindInventoryByProductId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateInventoryByProductId(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
