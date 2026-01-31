package model

import (
	"time"

	"codewithumam-kasir-api/internal/utils"
	"github.com/google/uuid"
)

// TODO: implement the metadata
type CategoryEntity struct {
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
	DeletedAt   *time.Time
	Version     int
	ID          uuid.UUID //UUIDv7
	Name        string
	Description string
}

type Category struct {
	ID          string     `json:"id"` //Base62 of UUIDv7
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int        `json:"version,omitempty"`
}

func (c *CategoryEntity) ToModel() *Category {
	return &Category{
		ID:          utils.EncodeBase62(c.ID.String()),
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   c.DeletedAt,
		Version:     c.Version,
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
		CreatedBy:   "USER",
		UpdatedBy:   "USER",
	}
}

// TODO: add validation
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     int    `json:"version"`
}

func (c *UpdateCategoryRequest) ToEntity() *CategoryEntity {
	return &CategoryEntity{
		Name:        c.Name,
		Description: c.Description,
		Version:     c.Version,
		UpdatedBy:   "USER",
	}
}
