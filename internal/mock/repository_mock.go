package mock

import (
	"codewithumam-kasir-api/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository is a mock implementation of CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) FindCategories() ([]model.CategoryEntity, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) FindCategoryByID(id string) (model.CategoryEntity, error) {
	args := m.Called(id)
	return args.Get(0).(model.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) FindCategoryByName(name string) (model.CategoryEntity, error) {
	args := m.Called(name)
	return args.Get(0).(model.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) InsertCategory(category model.CategoryEntity) (model.CategoryEntity, error) {
	args := m.Called(category)
	return args.Get(0).(model.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) UpdateCategoryByID(id string, category model.CategoryEntity) (model.CategoryEntity, error) {
	args := m.Called(id, category)
	return args.Get(0).(model.CategoryEntity), args.Error(1)
}

func (m *MockCategoryRepository) DeleteCategoryByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindProducts() ([]model.ProductEntity, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) FindProductByID(id string) (model.ProductEntity, error) {
	args := m.Called(id)
	return args.Get(0).(model.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) FindProductsByNameAndActiveStatus(name string, activeStatus *bool) ([]model.ProductEntity, error) {
	args := m.Called(name, activeStatus)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) InsertProduct(product model.ProductEntity) (model.ProductEntity, error) {
	args := m.Called(product)
	return args.Get(0).(model.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) UpdateProductByID(id string, product model.ProductEntity) (model.ProductEntity, error) {
	args := m.Called(id, product)
	return args.Get(0).(model.ProductEntity), args.Error(1)
}

func (m *MockProductRepository) DeleteProductByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
