package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/statistics", Statistics).Methods("GET")
	router.HandleFunc("/transaction", DeleteTransactions).Methods("DELETE")
	router.HandleFunc("/location", AddLocation).Methods("POST")
	router.HandleFunc("/location", UpdateLocation).Methods("PUT")

	return router
}

func TestStatistics(t *testing.T) {
	request, _ := http.NewRequest("GET", "/statistics", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestDeleteTransaction(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/transaction", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, 204, response.Code, "OK response is expected")
}

func TestAddLocation(t *testing.T) {
	location := &LocationDetails{
		Location: "bangalore",
	}
	jsonLocation, _ := json.Marshal(location)
	request, _ := http.NewRequest("POST", "/location", bytes.NewBuffer(jsonLocation))
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestUpdateLocation(t *testing.T) {
	location := &LocationDetails{
		Location: "bangalore",
	}
	jsonLocation, _ := json.Marshal(location)
	request, _ := http.NewRequest("PUT", "/location", bytes.NewBuffer(jsonLocation))
	response := httptest.NewRecorder()
	NewRouter().ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}
