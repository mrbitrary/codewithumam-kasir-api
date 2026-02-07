package repository

import (
	"testing"
	"time"

	"codewithumam-kasir-api/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryProductRepository_FindProducts(t *testing.T) {
	repo := NewProductRepository()

	// Initially empty
	products, err := repo.FindProducts()
	require.NoError(t, err)
	assert.Empty(t, products)

	// Add some products
	catID1 := uuid.New()
	prod1 := model.ProductEntity{
		ID:           uuid.New(),
		Name:         "Laptop",
		Price:        1000,
		Stocks:       5,
		CategoryID:   &catID1,
		CategoryName: "Electronics",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_, err = repo.InsertProduct(prod1)
	require.NoError(t, err)

	products, err = repo.FindProducts()
	require.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, "Laptop", products[0].Name)
}

func TestInMemoryProductRepository_FindProductByID(t *testing.T) {
	repo := NewProductRepository()

	id := uuid.New()
	catID2 := uuid.New()
	prod := model.ProductEntity{
		ID:           id,
		Name:         "Mouse",
		Price:        25,
		Stocks:       100,
		CategoryID:   &catID2,
		CategoryName: "Accessories",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_, err := repo.InsertProduct(prod)
	require.NoError(t, err)

	// Find existing product
	found, err := repo.FindProductByID(id.String())
	require.NoError(t, err)
	assert.Equal(t, "Mouse", found.Name)
	assert.Equal(t, int64(25), found.Price)

	// Find non-existent product
	_, err = repo.FindProductByID(uuid.New().String())
	assert.Error(t, err)
}

func TestInMemoryProductRepository_InsertProduct(t *testing.T) {
	repo := NewProductRepository()

	catID3 := uuid.New()
	prod := model.ProductEntity{
		ID:           uuid.New(),
		Name:         "Keyboard",
		Price:        75,
		Stocks:       50,
		CategoryID:   &catID3,
		CategoryName: "Accessories",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	inserted, err := repo.InsertProduct(prod)
	require.NoError(t, err)
	assert.Equal(t, prod.Name, inserted.Name)
	assert.Equal(t, prod.Price, inserted.Price)

	// Verify it's in the repository
	products, _ := repo.FindProducts()
	assert.Len(t, products, 1)
}

func TestInMemoryProductRepository_UpdateProductByID(t *testing.T) {
	repo := NewProductRepository()

	id := uuid.New()
	categoryID := uuid.New()
	prod := model.ProductEntity{
		ID:           id,
		Name:         "Original Product",
		Price:        100,
		Stocks:       10,
		CategoryID:   &categoryID,
		CategoryName: "Original Category",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_, err := repo.InsertProduct(prod)
	require.NoError(t, err)

	// Update product
	updated := model.ProductEntity{
		Name:         "Updated Product",
		Price:        150,
		Stocks:       20,
		CategoryID:   &categoryID,
		CategoryName: "Updated Category",
		UpdatedAt:    time.Now(),
	}
	result, err := repo.UpdateProductByID(id.String(), updated)
	require.NoError(t, err)
	assert.Equal(t, "Updated Product", result.Name)
	assert.Equal(t, int64(150), result.Price)
	assert.Equal(t, 20, result.Stocks)
	assert.Equal(t, id, result.ID)

	// Verify update persisted
	found, _ := repo.FindProductByID(id.String())
	assert.Equal(t, "Updated Product", found.Name)
	assert.Equal(t, int64(150), found.Price)

	// Update non-existent product
	_, err = repo.UpdateProductByID(uuid.New().String(), updated)
	assert.Error(t, err)
}

func TestInMemoryProductRepository_DeleteProductByID(t *testing.T) {
	repo := NewProductRepository()

	id := uuid.New()
	catID4 := uuid.New()
	prod := model.ProductEntity{
		ID:           id,
		Name:         "ToDelete",
		Price:        50,
		Stocks:       5,
		CategoryID:   &catID4,
		CategoryName: "Test",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_, err := repo.InsertProduct(prod)
	require.NoError(t, err)

	// Delete product
	err = repo.DeleteProductByID(id.String())
	require.NoError(t, err)

	// Verify deletion
	products, _ := repo.FindProducts()
	assert.Empty(t, products)

	// Delete non-existent product
	err = repo.DeleteProductByID(uuid.New().String())
	assert.Error(t, err)
}

func TestInMemoryProductRepository_InvalidUUID(t *testing.T) {
	repo := NewProductRepository()

	// FindProductByID with invalid UUID
	_, err := repo.FindProductByID("invalid-uuid")
	assert.Error(t, err)

	// UpdateProductByID with invalid UUID
	_, err = repo.UpdateProductByID("invalid-uuid", model.ProductEntity{})
	assert.Error(t, err)

	// DeleteProductByID with invalid UUID
	err = repo.DeleteProductByID("invalid-uuid")
	assert.Error(t, err)
}

func TestInMemoryProductRepository_MultipleProducts(t *testing.T) {
	repo := NewProductRepository()

	// Insert multiple products
	for i := 0; i < 5; i++ {
		catID5 := uuid.New()
		prod := model.ProductEntity{
			ID:           uuid.New(),
			Name:         "Product " + string(rune('A'+i)),
			Price:        int64((i + 1) * 100),
			Stocks:       (i + 1) * 10,
			CategoryID:   &catID5,
			CategoryName: "Category",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		_, err := repo.InsertProduct(prod)
		require.NoError(t, err)
	}

	// Verify all products are stored
	products, err := repo.FindProducts()
	require.NoError(t, err)
	assert.Len(t, products, 5)
}

func TestInMemoryProductRepository_FindProductsByNameAndActiveStatus(t *testing.T) {
	repo := NewProductRepository()

	id1 := uuid.New()
	id2 := uuid.New()
	deletedAt := time.Now()

	products := []model.ProductEntity{
		{ID: id1, Name: "Apple iPhone", DeletedAt: nil},
		{ID: id2, Name: "Samsung Galaxy", DeletedAt: &deletedAt},
	}

	for _, p := range products {
		_, _ = repo.InsertProduct(p)
	}

	// Test case: Active only
	active := true
	results, err := repo.FindProductsByNameAndActiveStatus("Apple", &active)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, id1, results[0].ID)

	// Test case: Inactive only
	inactive := false
	results, err = repo.FindProductsByNameAndActiveStatus("Samsung", &inactive)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, id2, results[0].ID)

	// Test case: All (nil activeStatus)
	results, err = repo.FindProductsByNameAndActiveStatus("", nil)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Test case: No results
	results, err = repo.FindProductsByNameAndActiveStatus("Nokia", nil)
	require.NoError(t, err)
	assert.Empty(t, results)
}
