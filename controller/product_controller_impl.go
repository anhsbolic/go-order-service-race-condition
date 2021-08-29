package controller

import (
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/anhsbolic/go-order-service-race-condition/model/web"
	"github.com/anhsbolic/go-order-service-race-condition/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type ProductControllerImpl struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &ProductControllerImpl{ProductService: productService}
}

func (controller *ProductControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productResponses := controller.ProductService.FindAll(request.Context())
	jsonResponse := web.JSONResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponses,
	}

	helper.WriteToResponseBody(writer, jsonResponse)
}

func (controller *ProductControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, err := strconv.Atoi(productId)
	helper.PanicIfError(err)

	productResponse := controller.ProductService.FindById(request.Context(), id)
	jsonResponse := web.JSONResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponse,
	}

	helper.WriteToResponseBody(writer, jsonResponse)
}

func (controller *ProductControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productCreateRequest := web.ProductCreateRequest{}
	helper.ReadFromRequestBody(request, &productCreateRequest)

	productResponse := controller.ProductService.Create(request.Context(), productCreateRequest)
	jsonResponse := web.JSONResponse{
		Code:   200,
		Status: "OK",
		Data:   productResponse,
	}

	helper.WriteToResponseBody(writer, jsonResponse)
}

func (controller *ProductControllerImpl) FindInventoryByProductId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, err := strconv.Atoi(productId)
	helper.PanicIfError(err)

	inventoryProductResponse := controller.ProductService.FindInventoryByProductId(request.Context(), id)
	jsonResponse := web.JSONResponse{
		Code:   200,
		Status: "OK",
		Data:   inventoryProductResponse,
	}

	helper.WriteToResponseBody(writer, jsonResponse)
}

func (controller *ProductControllerImpl) UpdateInventoryByProductId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	productId := params.ByName("productId")
	id, err := strconv.Atoi(productId)
	helper.PanicIfError(err)

	inventoryUpdateRequest := web.InventoryUpdateRequest{}
	helper.ReadFromRequestBody(request, &inventoryUpdateRequest)

	inventoryResponse := controller.ProductService.UpdateInventoryByProductId(request.Context(), id, inventoryUpdateRequest)
	jsonResponse := web.JSONResponse{
		Code:   200,
		Status: "OK",
		Data:   inventoryResponse,
	}

	helper.WriteToResponseBody(writer, jsonResponse)
}
