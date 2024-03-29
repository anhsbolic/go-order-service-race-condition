package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
	"github.com/anhsbolic/go-order-service-race-condition/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestGetProductsSuccess(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product1 := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Kaos Polos",
	})
	product2 := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Celana Panjang",
	})
	tx.Commit()

	router := SetupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	var products = responseBody["data"].([]interface{})
	productResponse1 := products[0].(map[string]interface{})
	productResponse2 := products[1].(map[string]interface{})

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	assert.Equal(t, product1.Id, int(productResponse1["id"].(float64)))
	assert.Equal(t, product1.Name, productResponse1["name"])

	assert.Equal(t, product2.Id, int(productResponse2["id"].(float64)))
	assert.Equal(t, product2.Name, productResponse2["name"])
}

func TestGetProductByIdSuccess(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Kaos Polos",
	})
	tx.Commit()

	router := SetupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products/"+strconv.Itoa(product.Id), nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, product.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, product.Name, responseBody["data"].(map[string]interface{})["name"])
}

func TestGetProductByIdFailed(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)
	router := SetupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products/10000", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestCreateProductSuccess(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)
	router := SetupRouter(db)

	requestBody := strings.NewReader(`{"name" : "Kaos Polos Hitam"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 201, response.StatusCode)

	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 201, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Kaos Polos Hitam", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateProductFailed(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)
	router := SetupRouter(db)

	requestBody := strings.NewReader(`{"name" : ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestGetInventoryByProductIdSuccess(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)

	inventoryRepository := repository.NewInventoryRepository()
	productRepository := repository.NewProductRepository()

	tx, _ := db.Begin()

	product := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Kaos Polos",
	})

	inventory := inventoryRepository.Save(context.Background(), tx, entity.Inventory{
		ProductId:      product.Id,
		StoredStock:    100,
		AvailableStock: 100,
		ReservedStock:  0,
	})

	tx.Commit()

	router := SetupRouter(db)
	url := fmt.Sprintf(`http://localhost:3000/api/products/%s/inventory`, strconv.Itoa(product.Id))
	request := httptest.NewRequest(http.MethodGet, url, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, product.Name, responseBody["data"].(map[string]interface{})["product_name"])
	assert.Equal(t, inventory.StoredStock, int(responseBody["data"].(map[string]interface{})["stored_stock"].(float64)))
	assert.Equal(t, inventory.AvailableStock, int(responseBody["data"].(map[string]interface{})["available_stock"].(float64)))
	assert.Equal(t, inventory.ReservedStock, int(responseBody["data"].(map[string]interface{})["reserved_stock"].(float64)))
}

func TestGetInventoryByProductIdFailed(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)

	router := SetupRouter(db)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/products/20000/inventory", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 404, response.StatusCode)
	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}
