package model

import (
	"time"

	"codewithumam-kasir-api/internal/utils"
	"github.com/google/uuid"
)

type TransactionEntity struct {
	ID                uuid.UUID
	TotalItems        int
	TotalPriceAmount  int64
	TotalPriceScale   int
	TotalPriceDisplay float64
	Currency          string
	CreatedAt         time.Time
	CreatedBy         string
	UpdatedAt         time.Time
	UpdatedBy         string
	DeletedAt         *time.Time
	Version           int
}

type TransactionDetailEntity struct {
	ID                uuid.UUID
	TransactionID     uuid.UUID
	ProductID         *uuid.UUID
	ProductName       string
	CategoryID        *uuid.UUID
	CategoryName      string
	PriceAmount       int64
	PriceScale        int
	PriceDisplay      float64
	Currency          string
	Quantity          int
	TotalPriceAmount  int64
	TotalPriceScale   int
	TotalPriceDisplay float64
	CreatedAt         time.Time
	CreatedBy         string
	UpdatedAt         time.Time
	UpdatedBy         string
	DeletedAt         *time.Time
	Version           int
}

type Transaction struct {
	ID         string              `json:"id"`
	TotalItems int                 `json:"total_items"`
	TotalPrice Price               `json:"total_price"`
	CreatedAt  time.Time           `json:"created_at"`
	Details    []TransactionDetail `json:"details,omitempty"`
}

type TransactionDetail struct {
	ProductID    string `json:"product_id,omitempty"`
	ProductName  string `json:"product_name"`
	CategoryID   string `json:"category_id,omitempty"`
	CategoryName string `json:"category_name"`
	Price        Price  `json:"price"`
	Quantity     int    `json:"quantity"`
	TotalPrice   Price  `json:"total_price"`
}

type Price struct {
	Amount   int64   `json:"amount"`
	Scale    int     `json:"scale"`
	Display  float64 `json:"display"`
	Currency string  `json:"currency"`
}

type CreateTransactionRequest struct {
	Items []CreateTransactionItemRequest `json:"items"`
}

type CreateTransactionItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (e *TransactionEntity) ToModel() *Transaction {
	return &Transaction{
		ID:         utils.EncodeBase62(e.ID.String()),
		TotalItems: e.TotalItems,
		TotalPrice: Price{
			Amount:   e.TotalPriceAmount,
			Scale:    e.TotalPriceScale,
			Display:  e.TotalPriceDisplay,
			Currency: e.Currency,
		},
		CreatedAt: e.CreatedAt,
	}
}

func (e *TransactionDetailEntity) ToModel() *TransactionDetail {
	var pID, cID string
	if e.ProductID != nil {
		pID = utils.EncodeBase62(e.ProductID.String())
	}
	if e.CategoryID != nil {
		cID = utils.EncodeBase62(e.CategoryID.String())
	}

	return &TransactionDetail{
		ProductID:    pID,
		ProductName:  e.ProductName,
		CategoryID:   cID,
		CategoryName: e.CategoryName,
		Price: Price{
			Amount:   e.PriceAmount,
			Scale:    e.PriceScale,
			Display:  e.PriceDisplay,
			Currency: e.Currency,
		},
		Quantity: e.Quantity,
		TotalPrice: Price{
			Amount:   e.TotalPriceAmount,
			Scale:    e.TotalPriceScale,
			Display:  e.TotalPriceDisplay,
			Currency: e.Currency,
		},
	}
}

type ReportResponse struct {
	TotalRevenue         Price             `json:"total_revenue"`
	TotalTransactions    int               `json:"total_transactions"`
	TopPopularItems      []PopularItem     `json:"top_popular_items"`
	TopPopularCategories []PopularCategory `json:"top_popular_categories"`
}

type PopularItem struct {
	Name         string `json:"name"`
	TotalSoldQty int    `json:"total_sold_qty"`
}

type PopularCategory struct {
	Name         string `json:"name"`
	TotalSoldQty int    `json:"total_sold_qty"`
}
