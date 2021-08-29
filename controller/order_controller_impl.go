package controller

import (
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/anhsbolic/go-order-service-race-condition/model/web"
	"github.com/anhsbolic/go-order-service-race-condition/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type OrderControllerImpl struct {
	OrderService service.OrderService
}

func NewOrderController(orderService service.OrderService) OrderController {
	return &OrderControllerImpl{OrderService: orderService}
}

func (controller *OrderControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	orderCreateRequest := web.OrderCreateRequest{}
	helper.ReadFromRequestBody(request, &orderCreateRequest)

	controller.OrderService.Create(request.Context(), orderCreateRequest)
	jsonResponse := web.JSONResponse{
		Code:   http.StatusCreated,
		Status: "OK",
	}

	writer.WriteHeader(http.StatusCreated)
	helper.WriteToResponseBody(writer, jsonResponse)
}
