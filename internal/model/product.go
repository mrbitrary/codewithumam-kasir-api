package model

import (
	"time"

	"codewithumam-kasir-api/internal/utils"
	"github.com/google/uuid"
)

// TODO: implement the metadata
// TODO: implement the category relationship
type ProductEntity struct {
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    *time.Time
	Version      int
	ID           uuid.UUID //UUIDv7
	Name         string
	Price        int64
	Stocks       int
	CategoryID   uuid.UUID
	CategoryName string // JOIN from category table by category_id
}

type Product struct {
	ID        string     `json:"id"` //Base62 of UUIDv7
	Name      string     `json:"name"`
	Price     int64      `json:"price"` //TODO: what is  the best way to represent price?
	Stocks    int        `json:"stocks"`
	Category  string     `json:"category,omitempty"` //category_name
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int        `json:"version,omitempty"`
}

func (p *ProductEntity) ToModel() *Product {
	return &Product{
		ID:        utils.EncodeBase62(p.ID.String()),
		Name:      p.Name,
		Price:     p.Price,
		Stocks:    p.Stocks,
		Category:  p.CategoryName,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: p.DeletedAt,
		Version:   p.Version,
	}
}

// TODO: add validation
type CreateProductRequest struct {
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Stocks   int    `json:"stocks"`
	Category string `json:"category"`
}

func (p *CreateProductRequest) ToEntity() *ProductEntity {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &ProductEntity{
		ID:           id,
		Name:         p.Name,
		Price:        p.Price,
		Stocks:       p.Stocks,
		CategoryName: p.Category,
		CreatedBy:    "USER",
		UpdatedBy:    "USER",
	}
}

// TODO: add validation
type UpdateProductRequest struct {
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Stocks   int    `json:"stocks"`
	Category string `json:"category"`
	Version  int    `json:"version"`
}

func (p *UpdateProductRequest) ToEntity() *ProductEntity {
	return &ProductEntity{
		Name:         p.Name,
		Price:        p.Price,
		Stocks:       p.Stocks,
		CategoryName: p.Category,
		Version:      p.Version,
		UpdatedBy:    "USER",
	}
}
