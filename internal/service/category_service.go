package service

import (
	"codewithumam-kasir-api/internal/models"
	"codewithumam-kasir-api/internal/repository"
	"codewithumam-kasir-api/internal/utils"
)

// TODO: optional try to implement partial update (PATCH)
// TODO: implement proper optimistic locking through etag
type CategoryService interface {
	FetchCategories() ([]models.Category, error)
	FetchCategoryByID(id string) (models.Category, error)
	CreateCategory(category models.CreateCategoryRequest) (models.Category, error)
	UpdateCategoryByID(id string, category models.UpdateCategoryRequest) (models.Category, error)
	DeleteCategoryByID(id string) error
}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{
		repository: repository,
	}
}

func (s *categoryService) FetchCategories() ([]models.Category, error) {
	entities, err := s.repository.FindCategories()
	if err != nil {
		return nil, err
	}
	categories := []models.Category{}
	for _, entity := range entities {
		categories = append(categories, *entity.ToModel())
	}
	return categories, nil
}

func (s *categoryService) FetchCategoryByID(id string) (models.Category, error) {
	entity, err := s.repository.FindCategoryByID(utils.DecodeBase62(id))
	if err != nil {
		return models.Category{}, err
	}
	return *entity.ToModel(), nil
}

func (s *categoryService) CreateCategory(request models.CreateCategoryRequest) (models.Category, error) {
	entity, err := s.repository.InsertCategory(*request.ToEntity())
	if err != nil {
		return models.Category{}, err
	}
	return *entity.ToModel(), nil
}

func (s *categoryService) UpdateCategoryByID(id string, request models.UpdateCategoryRequest) (models.Category, error) {
	entity, err := s.repository.UpdateCategoryByID(utils.DecodeBase62(id), *request.ToEntity())
	if err != nil {
		return models.Category{}, err
	}
	return *entity.ToModel(), nil
}

func (s *categoryService) DeleteCategoryByID(id string) error {
	return s.repository.DeleteCategoryByID(utils.DecodeBase62(id))
}
