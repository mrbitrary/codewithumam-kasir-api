package models

import (
	"codewithumam-kasir-api/internal/utils"
	"github.com/google/uuid"
)

// TODO: implement the metadata
type CategoryEntity struct {
	// CreatedAt   time.Time
	// CreatedBy	string
	// UpdatedAt   time.Time
	// UpdatedBy	string
	// DeletedAt   time.Time
	// Version     int
	ID          uuid.UUID //UUIDv7
	Name        string
	Description string
}

type Category struct {
	ID          string `json:"id"` //Base62 of UUIDv7
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *CategoryEntity) ToModel() *Category {
	return &Category{
		ID:          utils.EncodeBase62(c.ID.String()),
		Name:        c.Name,
		Description: c.Description,
	}
}

// TODO: add validation
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *CreateCategoryRequest) ToEntity() *CategoryEntity {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &CategoryEntity{
		ID:          id,
		Name:        c.Name,
		Description: c.Description,
	}
}

// TODO: add validation
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *UpdateCategoryRequest) ToEntity() *CategoryEntity {
	return &CategoryEntity{
		Name:        c.Name,
		Description: c.Description,
	}
}
