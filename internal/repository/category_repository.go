package repository

import (
	"codewithumam-kasir-api/internal/model"
)

type CategoryRepository interface {
	FindCategories() ([]model.CategoryEntity, error)
	FindCategoryByID(id string) (model.CategoryEntity, error)
	FindCategoryByName(name string) (model.CategoryEntity, error)
	InsertCategory(category model.CategoryEntity) (model.CategoryEntity, error)
	UpdateCategoryByID(id string, category model.CategoryEntity) (model.CategoryEntity, error)
	DeleteCategoryByID(id string) error
}
