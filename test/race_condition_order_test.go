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

func TestRaceConditionOrder(t *testing.T) {
	db := SetupTestDB()
	ClearProductsTable(db)
	ClearOrdersTable(db)

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

	i := 0
	for {
		total := 1
		bodyString := fmt.Sprintf(`{"product_id": %d, "total" : %d}`, product.Id, total)
		requestBody := strings.NewReader(bodyString)
		request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/orders", requestBody)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		response := recorder.Result()
		fmt.Println(fmt.Sprintf(`at order : %d `, i+1))
		fmt.Println(response)

		body, _ := ioutil.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		i++

		if i == inventory.AvailableStock+1 {
			assert.Equal(t, 400, response.StatusCode)
			assert.Equal(t, 400, int(responseBody["code"].(float64)))
			assert.Equal(t, "BAD REQUEST", responseBody["status"])

			break
		} else {
			assert.Equal(t, 201, response.StatusCode)
			assert.Equal(t, 201, int(responseBody["code"].(float64)))
			assert.Equal(t, "OK", responseBody["status"])
		}
	}
}
