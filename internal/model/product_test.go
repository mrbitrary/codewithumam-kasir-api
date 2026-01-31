package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductEntity_ToModel(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	categoryID := uuid.New()
	
	entity := &ProductEntity{
		ID:           id,
		Name:         "Laptop",
		Price:        150000,
		Stocks:       10,
		CategoryID:   categoryID,
		CategoryName: "Electronics",
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil,
		Version:      1,
	}

	model := entity.ToModel()

	require.NotNil(t, model)
	assert.NotEmpty(t, model.ID)
	assert.Equal(t, "Laptop", model.Name)
	assert.Equal(t, int64(150000), model.Price)
	assert.Equal(t, 10, model.Stocks)
	assert.Equal(t, "Electronics", model.Category)
	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, now, model.UpdatedAt)
	assert.Nil(t, model.DeletedAt)
	assert.Equal(t, 1, model.Version)
}

func TestProductEntity_ToModel_WithDeletedAt(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour)
	id := uuid.New()
	
	entity := &ProductEntity{
		ID:        id,
		Name:      "Discontinued",
		Price:     1000,
		Stocks:    0,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: &deletedAt,
		Version:   2,
	}

	model := entity.ToModel()

	require.NotNil(t, model)
	require.NotNil(t, model.DeletedAt)
	assert.Equal(t, deletedAt, *model.DeletedAt)
}

func TestCreateProductRequest_ToEntity(t *testing.T) {
	req := &CreateProductRequest{
		Name:     "Mouse",
		Price:    25000,
		Stocks:   50,
		Category: "Accessories",
	}

	entity := req.ToEntity()

	require.NotNil(t, entity)
	assert.NotEqual(t, uuid.Nil, entity.ID)
	assert.Equal(t, "Mouse", entity.Name)
	assert.Equal(t, int64(25000), entity.Price)
	assert.Equal(t, 50, entity.Stocks)
	assert.Equal(t, "Accessories", entity.CategoryName)
	assert.Equal(t, "USER", entity.CreatedBy)
	assert.Equal(t, "USER", entity.UpdatedBy)
}

func TestUpdateProductRequest_ToEntity(t *testing.T) {
	req := &UpdateProductRequest{
		Name:     "Updated Product",
		Price:    30000,
		Stocks:   100,
		Category: "New Category",
		Version:  5,
	}

	entity := req.ToEntity()

	require.NotNil(t, entity)
	assert.Equal(t, "Updated Product", entity.Name)
	assert.Equal(t, int64(30000), entity.Price)
	assert.Equal(t, 100, entity.Stocks)
	assert.Equal(t, "New Category", entity.CategoryName)
	assert.Equal(t, 5, entity.Version)
	assert.Equal(t, "USER", entity.UpdatedBy)
}
