package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error occured while initialising the database")
	}

	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
	id int NOT NULL AUTO_INCREMENT,
	name varchar(255) NOT NULL,
	quantity int,
	price float(10, 7),
	PRIMARY KEY (id)
	);`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
	log.Println("clearTable")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func checkStatuscode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatuscode(t, http.StatusOK, response.Code)
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"mouse", "quantity":1, "price":100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Contetnt-Type", "application/json")

	response := sendRequest(req)
	checkStatuscode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "mouse" {
		t.Errorf("Expected name: %v, Got: %v", "mouse", m["name"])
	}

	log.Printf("%T", m["quantity"])

	if m["quantity"] != 1.0 {
		t.Errorf("Expected quantity: %v, Got: %v", 1, m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
	// Clear table and add a product.
	clearTable()
	addProduct("connector", 10, 10)

	// Retrieve & Check if the product is added.
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatuscode(t, http.StatusOK, response.Code)

	// Delete the newly added product.
	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatuscode(t, http.StatusOK, response.Code)

	// Retrieve it again to check whether it exists or not.
	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatuscode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	// clear table and add a product
	clearTable()
	addProduct("connector", 10, 10)

	// Retrieve the product.
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)

	// Take the reponse in JSON format & puts it in the map[string]interface{} format.
	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	// Create our our new product & send the UPDATE request.
	var product = []byte(`{"name":"connector", "quantity":1, "price":10}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	// Take the reponse in JSON format & puts it in the map[string]interface{} format.
	response = sendRequest(req)
	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	// Check for values of the old & new product
	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected id: %v, Got: %v", newValue["id"], oldValue["id"])
	}

	if oldValue["name"] != newValue["name"] {
		t.Errorf("Expected id: %v, Got: %v", newValue["name"], oldValue["name"])
	}

	if oldValue["price"] != newValue["price"] {
		t.Errorf("Expected id: %v, Got: %v", newValue["price"], oldValue["price"])
	}

	if oldValue["quantity"] == newValue["quantity"] {
		t.Errorf("Expected id: %v, Got: %v", newValue["quantity"], oldValue["quantity"])
	}
}
