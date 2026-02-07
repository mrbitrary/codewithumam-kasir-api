package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTransactionEntity_ToModel(t *testing.T) {
	id, _ := uuid.NewV7()
	now := time.Now()
	entity := &TransactionEntity{
		ID:                id,
		TotalItems:        5,
		TotalPriceAmount:  5000,
		TotalPriceScale:   0,
		TotalPriceDisplay: 5000.0,
		Currency:          "IDR",
		CreatedAt:         now,
	}

	model := entity.ToModel()

	assert.NotNil(t, model)
	assert.Equal(t, 5, model.TotalItems)
	assert.Equal(t, int64(5000), model.TotalPrice.Amount)
}

func TestTransactionDetailEntity_ToModel(t *testing.T) {
	pID, _ := uuid.NewV7()
	cID, _ := uuid.NewV7()
	entity := &TransactionDetailEntity{
		ProductID:    &pID,
		ProductName:  "Product A",
		CategoryID:   &cID,
		CategoryName: "Category A",
		PriceAmount:  1000,
		Quantity:     2,
		Currency:     "IDR",
	}

	model := entity.ToModel()

	assert.NotNil(t, model)
	assert.Equal(t, "Product A", model.ProductName)
	assert.Equal(t, 2, model.Quantity)
}
