package repository

import (
	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/repository"

	"fmt"
	"github.com/google/uuid"
)

type CategoryRepositoryInMemoryImpl struct {
	categories []model.CategoryEntity
}

func NewCategoryRepository() repository.CategoryRepository {
	return &CategoryRepositoryInMemoryImpl{
		categories: []model.CategoryEntity{},
	}
}

func (r *CategoryRepositoryInMemoryImpl) FindCategories() ([]model.CategoryEntity, error) {
	return r.categories, nil
}

func (r *CategoryRepositoryInMemoryImpl) FindCategoryByID(id string) (model.CategoryEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return model.CategoryEntity{}, fmt.Errorf("category not found")
	}
	for _, category := range r.categories {
		if category.ID == parsedID {
			return category, nil
		}
	}
	return model.CategoryEntity{}, fmt.Errorf("category not found")
}

func (r *CategoryRepositoryInMemoryImpl) FindCategoryByName(name string) (model.CategoryEntity, error) {
	for _, category := range r.categories {
		if category.Name == name {
			return category, nil
		}
	}
	return model.CategoryEntity{}, fmt.Errorf("category not found")
}

func (r *CategoryRepositoryInMemoryImpl) InsertCategory(category model.CategoryEntity) (model.CategoryEntity, error) {
	r.categories = append(r.categories, category)
	return category, nil
}

func (r *CategoryRepositoryInMemoryImpl) UpdateCategoryByID(id string, category model.CategoryEntity) (model.CategoryEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return model.CategoryEntity{}, fmt.Errorf("category not found")
	}
	for i, c := range r.categories {
		if c.ID == parsedID {
			category.ID = parsedID
			r.categories[i] = category
			return category, nil
		}
	}
	return model.CategoryEntity{}, fmt.Errorf("category not found")
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
