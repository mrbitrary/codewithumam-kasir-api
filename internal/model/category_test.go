package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryEntity_ToModel(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	
	entity := &CategoryEntity{
		ID:          id,
		Name:        "Electronics",
		Description: "Electronic devices",
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   nil,
		Version:     1,
	}

	model := entity.ToModel()

	require.NotNil(t, model)
	assert.NotEmpty(t, model.ID)
	assert.Equal(t, "Electronics", model.Name)
	assert.Equal(t, "Electronic devices", model.Description)
	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, now, model.UpdatedAt)
	assert.Nil(t, model.DeletedAt)
	assert.Equal(t, 1, model.Version)
}

func TestCategoryEntity_ToModel_WithDeletedAt(t *testing.T) {
	now := time.Now()
	deletedAt := now.Add(time.Hour)
	id := uuid.New()
	
	entity := &CategoryEntity{
		ID:          id,
		Name:        "Archived",
		Description: "Archived category",
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   &deletedAt,
		Version:     2,
	}

	model := entity.ToModel()

	require.NotNil(t, model)
	require.NotNil(t, model.DeletedAt)
	assert.Equal(t, deletedAt, *model.DeletedAt)
}

func TestCreateCategoryRequest_ToEntity(t *testing.T) {
	req := &CreateCategoryRequest{
		Name:        "Books",
		Description: "Book collection",
	}

	entity := req.ToEntity()

	require.NotNil(t, entity)
	assert.NotEqual(t, uuid.Nil, entity.ID)
	assert.Equal(t, "Books", entity.Name)
	assert.Equal(t, "Book collection", entity.Description)
	assert.Equal(t, "USER", entity.CreatedBy)
	assert.Equal(t, "USER", entity.UpdatedBy)
}

func TestUpdateCategoryRequest_ToEntity(t *testing.T) {
	req := &UpdateCategoryRequest{
		Name:        "Updated Name",
		Description: "Updated Description",
		Version:     3,
	}

	entity := req.ToEntity()

	require.NotNil(t, entity)
	assert.Equal(t, "Updated Name", entity.Name)
	assert.Equal(t, "Updated Description", entity.Description)
	assert.Equal(t, 3, entity.Version)
	assert.Equal(t, "USER", entity.UpdatedBy)
}
