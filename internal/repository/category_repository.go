package repository

import (
	"codewithumam-kasir-api/internal/models"
)

type CategoryRepository interface {
	FindCategories() ([]models.CategoryEntity, error)
	FindCategoryByID(id string) (models.CategoryEntity, error)
	FindCategoryByName(name string) (models.CategoryEntity, error)
	InsertCategory(category models.CategoryEntity) (models.CategoryEntity, error)
	UpdateCategoryByID(id string, category models.CategoryEntity) (models.CategoryEntity, error)
	DeleteCategoryByID(id string) error
}
