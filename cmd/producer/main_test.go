package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/IBM/sarama/mocks"
	"github.com/gin-gonic/gin"
	"indasto1.com/unit3-mhs/broker"
	"indasto1.com/unit3-mhs/models"
)

var (
	router = CreateRouter()
)

func TestDataProducerWithValidData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	producer := mocks.NewSyncProducer(t, mocks.NewTestConfig())
	producer.ExpectSendMessageAndSucceed()
	broker.Producer = producer

	w := httptest.NewRecorder()

	exampleSum := models.Sum{
		Num1: 2,
		Num2: 3,
	}
	sumJson, _ := json.Marshal(exampleSum)
	req, _ := http.NewRequest("POST", "/data", strings.NewReader(string(sumJson)))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", w.Code, http.StatusOK)
	}

	expectedBody := ""
	if w.Body.String() != expectedBody {
		t.Errorf("handler returned wrong body: got '%v' expected '%v'", w.Body.String(), expectedBody)
	}
}

func TestDataProducerWithInvalidData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()

	exampleSum := map[string]string{
		"num1": "2",
		"num2": "3",
	}
	sumJson, _ := json.Marshal(exampleSum)
	req, _ := http.NewRequest("POST", "/data", strings.NewReader(string(sumJson)))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v expected %v", w.Code, http.StatusOK)
	}

	expectedBody := "invalid request body"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("handler returned wrong body: got '%v' expected '%v'", w.Body.String(), expectedBody)
	}
}
