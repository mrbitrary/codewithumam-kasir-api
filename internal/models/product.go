package models

import (
	"github.com/google/uuid"
	"kasir-api/internal/utils"
)

// TODO: implement the metadata
// TODO: implement the category relationship
type ProductEntity struct {
	// CreatedAt   time.Time
	// CreatedBy   string
	// UpdatedAt   time.Time
	// UpdatedBy   string
	// DeletedAt   time.Time
	// Version     int
	ID     uuid.UUID //UUIDv7
	Name   string
	Price  int64
	Stocks int
	// Category    string
}

type Product struct {
	ID     string `json:"id"` //Base62 of UUIDv7
	Name   string `json:"name"`
	Price  int64  `json:"price"` //TODO: what is  the best way to represent price?
	Stocks int    `json:"stocks"`
	// Category string `json:"category"`
}

func (p *ProductEntity) ToModel() *Product {
	return &Product{
		ID:     utils.EncodeBase62(p.ID.String()),
		Name:   p.Name,
		Price:  p.Price,
		Stocks: p.Stocks,
		// Category:    p.Category,
	}
}

// TODO: add validation
type CreateProductRequest struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Stocks int    `json:"stocks"`
	// Category    string `json:"category"`
}

func (p *CreateProductRequest) ToEntity() *ProductEntity {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &ProductEntity{
		ID:     id,
		Name:   p.Name,
		Price:  p.Price,
		Stocks: p.Stocks,
		// Category:    p.Category,
	}
}

// TODO: add validation
type UpdateProductRequest struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Stocks int    `json:"stocks"`
	// Category    string `json:"category"`
}

func (p *UpdateProductRequest) ToEntity() *ProductEntity {
	return &ProductEntity{
		Name:   p.Name,
		Price:  p.Price,
		Stocks: p.Stocks,
		// Category:    p.Category,
	}
}
