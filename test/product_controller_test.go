package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/anhsbolic/go-order-service-race-condition/app"
	"github.com/anhsbolic/go-order-service-race-condition/helper"
	"github.com/anhsbolic/go-order-service-race-condition/model/entity"
	"github.com/anhsbolic/go-order-service-race-condition/repository"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:newpass@tcp(localhost:3306)/order-service-race-condition-poc")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	router := app.NewRouter(db, validate)

	return router
}

func truncateProductsTable(db *sql.DB) {
	db.Exec("truncate products")
}

func TestGetProductsSuccess(t *testing.T) {
	db := setupTestDB()
	truncateProductsTable(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product1 := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Kaos Polos",
	})
	product2 := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Celana Panjang",
	})
	tx.Commit()

	router := setupRouter(db)
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
	db := setupTestDB()
	truncateProductsTable(db)

	tx, _ := db.Begin()
	productRepository := repository.NewProductRepository()
	product := productRepository.Save(context.Background(), tx, entity.Product{
		Name: "Kaso Polos",
	})
	tx.Commit()

	router := setupRouter(db)
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
	db := setupTestDB()
	truncateProductsTable(db)
	router := setupRouter(db)

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
	db := setupTestDB()
	truncateProductsTable(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name" : "Kaos Polos Hitam"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/products", requestBody)
	request.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := ioutil.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, "Kaos Polos Hitam", responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateProductFailed(t *testing.T) {
	db := setupTestDB()
	truncateProductsTable(db)
	router := setupRouter(db)

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
