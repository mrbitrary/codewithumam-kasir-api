package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "codewithumam-kasir-api/internal/mock"
	"codewithumam-kasir-api/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductHandlerFetchProducts(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	products := []model.Product{
		{ID: "1", Name: "Laptop", Price: 1000, Stocks: 5, Category: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
		{ID: "2", Name: "Mouse", Price: 25, Stocks: 50, Category: "Accessories", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
	}

	mockService.On("FetchProducts").Return(products, nil)

	req := httptest.NewRequest("GET", "/api/products", nil)
	rec := httptest.NewRecorder()

	handler.FetchProducts(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.APIResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response.Data)
	mockService.AssertExpectations(t)
}

func TestProductHandlerFetchProducts_WithFilters(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	products := []model.Product{
		{ID: "1", Name: "Laptop", Price: 1000},
	}

	active := true
	mockService.On("FetchProductsByNameAndActiveStatus", "Laptop", &active).Return(products, nil)

	req := httptest.NewRequest("GET", "/api/products?name=Laptop&active=1", nil)
	rec := httptest.NewRecorder()

	handler.FetchProducts(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandlerFetchProducts_InvalidActive(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest("GET", "/api/products?active=invalid", nil)
	rec := httptest.NewRecorder()

	handler.FetchProducts(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandlerFetchProductsError(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	mockService.On("FetchProducts").Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/api/products", nil)
	rec := httptest.NewRecorder()

	handler.FetchProducts(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandlerFetchProductByID(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	product := model.Product{
		ID: "1", Name: "Laptop", Price: 1000, Stocks: 5,
		Category: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1,
	}

	mockService.On("FetchProductByID", "test-id").Return(product, nil)

	req := httptest.NewRequest("GET", "/api/products/test-id", nil)
	req.SetPathValue("id", "test-id")
	rec := httptest.NewRecorder()

	handler.FetchProductByID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandlerCreateProduct(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	reqBody := model.CreateProductRequest{
		Name:     "New Product",
		Price:    500,
		Stocks:   10,
		Category: "Electronics",
	}

	product := model.Product{
		ID: "1", Name: "New Product", Price: 500, Stocks: 10,
		Category: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1,
	}

	mockService.On("CreateProduct", reqBody).Return(product, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateProduct(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandlerCreateProductInvalidJSON(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest("POST", "/api/products", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateProduct(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProductHandlerUpdateProduct(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	reqBody := model.UpdateProductRequest{
		Name:     "Updated Product",
		Price:    600,
		Stocks:   20,
		Category: "Updated Category",
		Version:  2,
	}

	product := model.Product{
		ID: "1", Name: "Updated Product", Price: 600, Stocks: 20,
		Category: "Updated Category", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 2,
	}

	mockService.On("UpdateProductByID", "test-id", reqBody).Return(product, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/products/test-id", bytes.NewBuffer(body))
	req.SetPathValue("id", "test-id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateProduct(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandlerDeleteProduct(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	mockService.On("DeleteProductByID", "test-id").Return(nil)

	req := httptest.NewRequest("DELETE", "/api/products/test-id", nil)
	req.SetPathValue("id", "test-id")
	rec := httptest.NewRecorder()

	handler.DeleteProduct(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandlerDeleteProductError(t *testing.T) {
	mockService := new(mocks.MockProductService)
	handler := NewProductHandler(mockService)

	mockService.On("DeleteProductByID", "test-id").Return(errors.New("delete error"))

	req := httptest.NewRequest("DELETE", "/api/products/test-id", nil)
	req.SetPathValue("id", "test-id")
	rec := httptest.NewRecorder()

	handler.DeleteProduct(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}
