package repository

import (
	"codewithumam-kasir-api/internal/models"
	"fmt"
	"github.com/google/uuid"
)

type CategoryRepository interface {
	FindCategories() ([]models.CategoryEntity, error)
	FindCategoryByID(id string) (models.CategoryEntity, error)
	FindCategoryByName(name string) (models.CategoryEntity, error)
	InsertCategory(category models.CategoryEntity) (models.CategoryEntity, error)
	UpdateCategoryByID(id string, category models.CategoryEntity) (models.CategoryEntity, error)
	DeleteCategoryByID(id string) error
}

type CategoryRepositoryInMemoryImpl struct {
	categories []models.CategoryEntity
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryInMemoryImpl{
		categories: []models.CategoryEntity{},
	}
}

func (r *CategoryRepositoryInMemoryImpl) FindCategories() ([]models.CategoryEntity, error) {
	return r.categories, nil
}

func (r *CategoryRepositoryInMemoryImpl) FindCategoryByID(id string) (models.CategoryEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.CategoryEntity{}, fmt.Errorf("category not found")
	}
	for _, category := range r.categories {
		if category.ID == parsedID {
			return category, nil
		}
	}
	return models.CategoryEntity{}, fmt.Errorf("category not found")
}

func (r *CategoryRepositoryInMemoryImpl) FindCategoryByName(name string) (models.CategoryEntity, error) {
	for _, category := range r.categories {
		if category.Name == name {
			return category, nil
		}
	}
	return models.CategoryEntity{}, fmt.Errorf("category not found")
}

func (r *CategoryRepositoryInMemoryImpl) InsertCategory(category models.CategoryEntity) (models.CategoryEntity, error) {
	r.categories = append(r.categories, category)
	return category, nil
}

func (r *CategoryRepositoryInMemoryImpl) UpdateCategoryByID(id string, category models.CategoryEntity) (models.CategoryEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.CategoryEntity{}, fmt.Errorf("category not found")
	}
	for i, c := range r.categories {
		if c.ID == parsedID {
			category.ID = parsedID
			r.categories[i] = category
			return category, nil
		}
	}
	return models.CategoryEntity{}, fmt.Errorf("category not found")
}

func (r *CategoryRepositoryInMemoryImpl) DeleteCategoryByID(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("category not found")
	}
	for i, c := range r.categories {
		if c.ID == parsedID {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("category not found")
}
