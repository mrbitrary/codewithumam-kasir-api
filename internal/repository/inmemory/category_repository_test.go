package repository

import (
	"testing"
	"time"

	"codewithumam-kasir-api/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryCategoryRepository_FindCategories(t *testing.T) {
	repo := NewCategoryRepository()

	// Initially empty
	categories, err := repo.FindCategories()
	require.NoError(t, err)
	assert.Empty(t, categories)

	// Add some categories
	cat1 := model.CategoryEntity{
		ID:          uuid.New(),
		Name:        "Electronics",
		Description: "Electronic devices",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err = repo.InsertCategory(cat1)
	require.NoError(t, err)

	categories, err = repo.FindCategories()
	require.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, "Electronics", categories[0].Name)
}

func TestInMemoryCategoryRepository_FindCategoryByID(t *testing.T) {
	repo := NewCategoryRepository()

	id := uuid.New()
	cat := model.CategoryEntity{
		ID:          id,
		Name:        "Books",
		Description: "Reading materials",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := repo.InsertCategory(cat)
	require.NoError(t, err)

	// Find existing category
	found, err := repo.FindCategoryByID(id.String())
	require.NoError(t, err)
	assert.Equal(t, "Books", found.Name)

	// Find non-existent category
	_, err = repo.FindCategoryByID(uuid.New().String())
	assert.Error(t, err)
}

func TestInMemoryCategoryRepository_FindCategoryByName(t *testing.T) {
	repo := NewCategoryRepository()

	cat := model.CategoryEntity{
		ID:          uuid.New(),
		Name:        "Sports",
		Description: "Sports equipment",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := repo.InsertCategory(cat)
	require.NoError(t, err)

	// Find existing category
	found, err := repo.FindCategoryByName("Sports")
	require.NoError(t, err)
	assert.Equal(t, "Sports", found.Name)

	// Find non-existent category
	_, err = repo.FindCategoryByName("NonExistent")
	assert.Error(t, err)
}

func TestInMemoryCategoryRepository_InsertCategory(t *testing.T) {
	repo := NewCategoryRepository()

	cat := model.CategoryEntity{
		ID:          uuid.New(),
		Name:        "Furniture",
		Description: "Home furniture",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	inserted, err := repo.InsertCategory(cat)
	require.NoError(t, err)
	assert.Equal(t, cat.Name, inserted.Name)

	// Verify it's in the repository
	categories, _ := repo.FindCategories()
	assert.Len(t, categories, 1)
}

func TestInMemoryCategoryRepository_UpdateCategoryByID(t *testing.T) {
	repo := NewCategoryRepository()

	id := uuid.New()
	cat := model.CategoryEntity{
		ID:          id,
		Name:        "Original",
		Description: "Original description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := repo.InsertCategory(cat)
	require.NoError(t, err)

	// Update category
	updated := model.CategoryEntity{
		Name:        "Updated",
		Description: "Updated description",
		UpdatedAt:   time.Now(),
	}
	result, err := repo.UpdateCategoryByID(id.String(), updated)
	require.NoError(t, err)
	assert.Equal(t, "Updated", result.Name)
	assert.Equal(t, id, result.ID)

	// Verify update persisted
	found, _ := repo.FindCategoryByID(id.String())
	assert.Equal(t, "Updated", found.Name)

	// Update non-existent category
	_, err = repo.UpdateCategoryByID(uuid.New().String(), updated)
	assert.Error(t, err)
}

func TestInMemoryCategoryRepository_DeleteCategoryByID(t *testing.T) {
	repo := NewCategoryRepository()

	id := uuid.New()
	cat := model.CategoryEntity{
		ID:          id,
		Name:        "ToDelete",
		Description: "Will be deleted",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := repo.InsertCategory(cat)
	require.NoError(t, err)

	// Delete category
	err = repo.DeleteCategoryByID(id.String())
	require.NoError(t, err)

	// Verify deletion
	categories, _ := repo.FindCategories()
	assert.Empty(t, categories)

	// Delete non-existent category
	err = repo.DeleteCategoryByID(uuid.New().String())
	assert.Error(t, err)
}

func TestInMemoryCategoryRepository_InvalidUUID(t *testing.T) {
	repo := NewCategoryRepository()

	// FindCategoryByID with invalid UUID
	_, err := repo.FindCategoryByID("invalid-uuid")
	assert.Error(t, err)

	// UpdateCategoryByID with invalid UUID
	_, err = repo.UpdateCategoryByID("invalid-uuid", model.CategoryEntity{})
	assert.Error(t, err)

	// DeleteCategoryByID with invalid UUID
	err = repo.DeleteCategoryByID("invalid-uuid")
	assert.Error(t, err)
}
