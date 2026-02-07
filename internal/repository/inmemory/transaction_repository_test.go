package repository

import (
	"codewithumam-kasir-api/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransactionRepositoryInMemory_CreateTransaction(t *testing.T) {
	productRepo := NewProductRepository()
	txRepo := NewTransactionRepository(productRepo)

	productID, _ := uuid.NewV7()
	product := model.ProductEntity{
		ID:     productID,
		Name:   "Test Product",
		Price:  1000,
		Stocks: 10,
	}
	_, _ = productRepo.InsertProduct(product)

	txID, _ := uuid.NewV7()
	tx := model.TransactionEntity{
		ID:               txID,
		TotalItems:       2,
		TotalPriceAmount: 2000,
	}

	detailID, _ := uuid.NewV7()
	details := []model.TransactionDetailEntity{
		{
			ID:            detailID,
			TransactionID: txID,
			ProductID:     &productID,
			Quantity:      2,
		},
	}

	result, err := txRepo.CreateTransaction(tx, details)

	assert.NoError(t, err)
	assert.Equal(t, txID, result.ID)

	// Verify stock deduction
	updatedProduct, _ := productRepo.FindProductByID(productID.String())
	assert.Equal(t, 8, updatedProduct.Stocks)
}

func TestTransactionRepositoryInMemory_GetReports(t *testing.T) {
	productRepo := NewProductRepository()
	txRepo := NewTransactionRepository(productRepo)

	// Currently returns empty/nil as implemented
	resp, err := txRepo.GetReportStats(time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.TotalRevenue.Amount)

	cat, err := txRepo.GetMostPopularCategory(time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Empty(t, cat.Name)

	prod, err := txRepo.GetMostPopularProduct(time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Empty(t, prod.Name)
}
