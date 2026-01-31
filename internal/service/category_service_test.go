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

func TestCategoryServiceFetchCategories(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	now := time.Now()
	entities := []model.CategoryEntity{
		{ID: uuid.New(), Name: "Electronics", Description: "Devices", CreatedAt: now, UpdatedAt: now, Version: 1},
		{ID: uuid.New(), Name: "Books", Description: "Reading", CreatedAt: now, UpdatedAt: now, Version: 1},
	}

	mockRepo.On("FindCategories").Return(entities, nil)

	categories, err := service.FetchCategories()

	require.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Electronics", categories[0].Name)
	assert.Equal(t, "Books", categories[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryServiceFetchCategoriesError(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	mockRepo.On("FindCategories").Return(nil, errors.New("database error"))

	categories, err := service.FetchCategories()

	assert.Error(t, err)
	assert.Nil(t, categories)
	mockRepo.AssertExpectations(t)
}

func TestCategoryServiceFetchCategoryByID(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	now := time.Now()
	id := uuid.New()
	entity := model.CategoryEntity{
		ID: id, Name: "Electronics", Description: "Devices",
		CreatedAt: now, UpdatedAt: now, Version: 1,
	}

	mockRepo.On("FindCategoryByID", mock.Anything).Return(entity, nil)

	category, err := service.FetchCategoryByID("test-id")

	require.NoError(t, err)
	assert.Equal(t, "Electronics", category.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryServiceCreateCategory(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	request := model.CreateCategoryRequest{
		Name:        "New Category",
		Description: "New Description",
	}

	mockRepo.On("InsertCategory", mock.AnythingOfType("model.CategoryEntity")).
		Return(model.CategoryEntity{
			ID: uuid.New(), Name: "New Category", Description: "New Description",
			CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 1,
		}, nil)

	category, err := service.CreateCategory(request)

	require.NoError(t, err)
	assert.Equal(t, "New Category", category.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryServiceUpdateCategoryByID(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	request := model.UpdateCategoryRequest{
		Name:        "Updated",
		Description: "Updated Desc",
		Version:     2,
	}

	mockRepo.On("UpdateCategoryByID", mock.Anything, mock.AnythingOfType("model.CategoryEntity")).
		Return(model.CategoryEntity{
			ID: uuid.New(), Name: "Updated", Description: "Updated Desc",
			CreatedAt: time.Now(), UpdatedAt: time.Now(), Version: 2,
		}, nil)

	category, err := service.UpdateCategoryByID("test-id", request)

	require.NoError(t, err)
	assert.Equal(t, "Updated", category.Name)
	mockRepo.AssertExpectations(t)
}

func TestCategoryServiceDeleteCategoryByID(t *testing.T) {
	mockRepo := new(mocks.MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	mockRepo.On("DeleteCategoryByID", mock.Anything).Return(nil)

	err := service.DeleteCategoryByID("test-id")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
