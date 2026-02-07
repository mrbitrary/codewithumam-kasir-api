package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codewithumam-kasir-api/internal/mock"
	"codewithumam-kasir-api/internal/model"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	mockService := new(mock.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	reqBody := model.CreateTransactionRequest{
		Items: []model.CreateTransactionItemRequest{
			{ProductID: "abc", Quantity: 1},
		},
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	mockService.On("CreateTransaction", reqBody).Return(model.Transaction{ID: "tx_123"}, nil)

	handler.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response model.APIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestTransactionHandler_FetchReport(t *testing.T) {
	mockService := new(mock.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	tests := []struct {
		url    string
		period string
	}{
		{"/api/reports/today", "today"},
		{"/api/reports/yesterday", "yesterday"},
		{"/api/reports/last-week", "last-week"},
		{"/api/reports/last-month", "last-month"},
		{"/api/reports/week-to-date", "week-to-date"},
		{"/api/reports/month-to-date", "month-to-date"},
		{"/api/reports/year-to-date", "year-to-date"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()

			mockService.On("FetchReport", "", "", tt.period).Return(model.ReportResponse{TotalTransactions: 10}, nil)

			handler.FetchReport(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var response model.APIResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.NotNil(t, response.Data)
		})
	}
}

func TestTransactionHandler_FetchPopularCategory(t *testing.T) {
	mockService := new(mock.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest("GET", "/api/reports/popular-category", nil)
	rr := httptest.NewRecorder()

	mockService.On("FetchMostPopularCategory", "", "").Return(model.PopularCategory{Name: "Cat A", TotalSoldQty: 50}, nil)

	handler.FetchPopularCategory(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTransactionHandler_FetchReport_InvalidDateRange(t *testing.T) {
	mockService := new(mock.MockTransactionService)
	handler := NewTransactionHandler(mockService)

	req, _ := http.NewRequest("GET", "/api/reports?startDate=2024-01-02&endDate=2024-01-01", nil)
	rr := httptest.NewRecorder()

	mockService.On("FetchReport", "2024-01-02", "2024-01-01", "").Return(model.ReportResponse{}, errors.New("startDate cannot be after endDate"))

	handler.FetchReport(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
