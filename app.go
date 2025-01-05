package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize() error {
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

// General function to send any JSON response with a specified status code and payload.
func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)               // Converts the payload(any Go data structures) into a JSON
	w.Header().Set("content-type", "application/json") // Sets the content type to application/json
	w.WriteHeader(statusCode)                          // Sets the HTTP status code
	w.Write(response)                                  // Writes the JSON response to the response writer
}

// Convenience function to send an error message as a JSON response, utilizing the sendResponse function.
func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	// fetch the product id
	vars := mux.Vars(r) // localhost:10000/product/1
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	p := product{ID: key}
	err = p.getProduct(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows: // user asked for an id not in our DB.
			sendError(w, http.StatusNotFound, "Product not found ")
		default: // any other error
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product

	err := json.NewDecoder(r.Body).Decode(&p) // takes the json request and turns it into a go struct format (in our case is the product struct) & then put it in the p variable.
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = p.createProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	// fetch the product id
	vars := mux.Vars(r) // localhost:10000/product/1
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	// Take input and turn it from json to struct
	var p product
	err = json.NewDecoder(r.Body).Decode(&p) // takes the json request and turns it into a go struct format (in our case is the product struct) & then put it in the p variable.
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	//
	p.ID = key
	err = p.updateProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	// fetch the product id
	vars := mux.Vars(r) // localhost:10000/product/1
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	p := product{ID: key}
	err = p.deleteProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
	}

	sendResponse(w, http.StatusOK, map[string]string{"result": "successful deletion"})
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}
