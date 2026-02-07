package mock

import (
	"codewithumam-kasir-api/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockCategoryService is a mock implementation of CategoryService
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) FetchCategories() ([]model.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryService) FetchCategoryByID(id string) (model.Category, error) {
	args := m.Called(id)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *MockCategoryService) CreateCategory(category model.CreateCategoryRequest) (model.Category, error) {
	args := m.Called(category)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *MockCategoryService) UpdateCategoryByID(id string, category model.UpdateCategoryRequest) (model.Category, error) {
	args := m.Called(id, category)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *MockCategoryService) DeleteCategoryByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockProductService is a mock implementation of ProductService
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) FetchProducts() ([]model.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductService) FetchProductByID(id string) (model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockProductService) FetchProductsByNameAndActiveStatus(name string, activeStatus *bool) ([]model.Product, error) {
	args := m.Called(name, activeStatus)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductService) CreateProduct(product model.CreateProductRequest) (model.Product, error) {
	args := m.Called(product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockProductService) UpdateProductByID(id string, product model.UpdateProductRequest) (model.Product, error) {
	args := m.Called(id, product)
	return args.Get(0).(model.Product), args.Error(1)
}

func (m *MockProductService) DeleteProductByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockTransactionService is a mock implementation of TransactionService
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(req model.CreateTransactionRequest) (model.Transaction, error) {
	args := m.Called(req)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func (m *MockTransactionService) FetchReport(startDateStr, endDateStr, period string) (model.ReportResponse, error) {
	args := m.Called(startDateStr, endDateStr, period)
	return args.Get(0).(model.ReportResponse), args.Error(1)
}

func (m *MockTransactionService) FetchMostPopularCategory(startDateStr, endDateStr string) (model.PopularCategory, error) {
	args := m.Called(startDateStr, endDateStr)
	return args.Get(0).(model.PopularCategory), args.Error(1)
}

func (m *MockTransactionService) FetchMostPopularProduct(startDateStr, endDateStr string) (model.PopularItem, error) {
	args := m.Called(startDateStr, endDateStr)
	return args.Get(0).(model.PopularItem), args.Error(1)
}
