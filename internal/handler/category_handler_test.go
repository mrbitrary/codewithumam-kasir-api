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

func TestCategoryHandlerFetchCategories(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	categories := []model.Category{
		{ID: "1", Name: "Electronics", Description: "Devices", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
		{ID: "2", Name: "Books", Description: "Reading", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1},
	}

	mockService.On("FetchCategories").Return(categories, nil)

	req := httptest.NewRequest("GET", "/api/categories", nil)
	rec := httptest.NewRecorder()

	handler.FetchCategories(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.APIResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	require.NoError(t, err)
	assert.NotNil(t, response.Data)
	mockService.AssertExpectations(t)
}

func TestCategoryHandlerFetchCategoriesError(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	mockService.On("FetchCategories").Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/api/categories", nil)
	rec := httptest.NewRecorder()

	handler.FetchCategories(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandlerFetchCategoryByID(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	category := model.Category{
		ID: "1", Name: "Electronics", Description: "Devices",
		CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1,
	}

	mockService.On("FetchCategoryByID", "test-id").Return(category, nil)

	req := httptest.NewRequest("GET", "/api/categories/test-id", nil)
	req.SetPathValue("id", "test-id")
	rec := httptest.NewRecorder()

	handler.FetchCategoryByID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandlerCreateCategory(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	reqBody := model.CreateCategoryRequest{
		Name:        "New Category",
		Description: "New Description",
	}

	category := model.Category{
		ID: "1", Name: "New Category", Description: "New Description",
		CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1,
	}

	mockService.On("CreateCategory", reqBody).Return(category, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateCategory(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandlerCreateCategoryInvalidJSON(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	req := httptest.NewRequest("POST", "/api/categories", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateCategory(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCategoryHandlerUpdateCategory(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	reqBody := model.UpdateCategoryRequest{
		Name:        "Updated",
		Description: "Updated Desc",
		Version:     2,
	}

	category := model.Category{
		ID: "1", Name: "Updated", Description: "Updated Desc",
		CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 2,
	}

	mockService.On("UpdateCategoryByID", "test-id", reqBody).Return(category, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/api/categories/test-id", bytes.NewBuffer(body))
	req.SetPathValue("id", "test-id")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateCategory(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandlerDeleteCategory(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	mockService.On("DeleteCategoryByID", "test-id").Return(nil)

	req := httptest.NewRequest("DELETE", "/api/categories/test-id", nil)
	req.SetPathValue("id", "test-id")
	rec := httptest.NewRecorder()

	handler.DeleteCategory(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCategoryHandlerDeleteCategoryError(t *testing.T) {
	mockService := new(mocks.MockCategoryService)
	handler := NewCategoryHandler(mockService)

	mockService.On("DeleteCategoryByID", "test-id").Return(errors.New("delete error"))

	req := httptest.NewRequest("DELETE", "/api/categories/test-id", nil)
	req.SetPathValue("id", "test-id")
	rec := httptest.NewRecorder()

	handler.DeleteCategory(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}
