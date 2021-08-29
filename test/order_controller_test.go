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
	"strings"
	"testing"
)

func TestCreateOrderSuccess(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)
	ClearOrdersTable(db)

	inventoryRepository := repository.NewInventoryRepository()
	productRepository := repository.NewProductRepository()

	tx, _ := db.Begin()

	product := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Kaos Polos",
	})

	inventoryRepository.Save(context.Background(), tx, entity.Inventory{
		ProductId:      product.Id,
		StoredStock:    100,
		AvailableStock: 100,
		ReservedStock:  0,
	})

	tx.Commit()

	router := SetupRouter(db)
	total := 10
	bodyString := fmt.Sprintf(`{"product_id": %d, "total" : %d}`, product.Id, total)
	requestBody := strings.NewReader(bodyString)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/orders", requestBody)
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
}

func TestCreateOrderFailed(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)
	router := SetupRouter(db)

	requestBody := strings.NewReader(`{"product_id" : 0, "total" : 0}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/orders", requestBody)
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
