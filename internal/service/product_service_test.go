package service

import (
	"errors"
	"testing"
	"time"

	mocks "codewithumam-kasir-api/internal/mock"
	"codewithumam-kasir-api/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProductServiceFetchProducts(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	now := time.Now()
	entities := []model.ProductEntity{
		{ID: uuid.New(), Name: "Laptop", Price: 1000, Stocks: 5, CategoryName: "Electronics", CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: uuid.New(), Name: "Mouse", Price: 25, Stocks: 50, CategoryName: "Accessories", CreatedAt: now, UpdatedAt: now, Version: 1},
	}

	mockRepo.On("FindProducts").Return(entities, nil)

	products, err := service.FetchProducts()

	require.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Laptop", products[0].Name)
	assert.Equal(t, "Mouse", products[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestProductServiceFetchProductsError(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("FindProducts").Return(nil, errors.New("database error"))

	products, err := service.FetchProducts()

	assert.Error(t, err)
	assert.Nil(t, products)
	mockRepo.AssertExpectations(t)
}

func TestProductServiceFetchProductByID(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	now := time.Now()
	id := uuid.New()
	entity := model.ProductEntity{
		ID: id, Name: "Laptop", Price: 1000, Stocks: 5,
		CategoryName: "Electronics", CreatedAt: now, UpdatedAt: now, Version: 1,
	}

	mockRepo.On("FindProductByID", mock.Anything).Return(entity, nil)

	product, err := service.FetchProductByID("test-id")

	require.NoError(t, err)
	assert.Equal(t, "Laptop", product.Name)
	assert.Equal(t, int64(1000), product.Price)
	mockRepo.AssertExpectations(t)
}

func TestProductServiceCreateProduct(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	request := model.CreateProductRequest{
		Name:     "New Product",
		Price:    500,
		Stocks:   10,
		Category: "Electronics",
	}

	mockRepo.On("InsertProduct", mock.AnythingOfType("model.ProductEntity")).
		Return(model.ProductEntity{
			ID: uuid.New(), Name: "New Product", Price: 500, Stocks: 10,
			CategoryName: "Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1,
		}, nil)

	product, err := service.CreateProduct(request)

	require.NoError(t, err)
	assert.Equal(t, "New Product", product.Name)
	assert.Equal(t, int64(500), product.Price)
	mockRepo.AssertExpectations(t)
}

func TestProductServiceUpdateProductByID(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	request := model.UpdateProductRequest{
		Name:     "Updated Product",
		Price:    600,
		Stocks:   20,
		Category: "Updated Category",
		Version:  2,
	}

	mockRepo.On("UpdateProductByID", mock.Anything, mock.AnythingOfType("model.ProductEntity")).
		Return(model.ProductEntity{
			ID: uuid.New(), Name: "Updated Product", Price: 600, Stocks: 20,
			CategoryName: "Updated Category", CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 2,
		}, nil)

	product, err := service.UpdateProductByID("test-id", request)

	require.NoError(t, err)
	assert.Equal(t, "Updated Product", product.Name)
	assert.Equal(t, int64(600), product.Price)
	mockRepo.AssertExpectations(t)
}

func TestProductServiceDeleteProductByID(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("DeleteProductByID", mock.Anything).Return(nil)

	err := service.DeleteProductByID("test-id")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
