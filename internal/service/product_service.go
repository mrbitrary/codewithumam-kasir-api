package service

import (
	"codewithumam-kasir-api/internal/models"
	"codewithumam-kasir-api/internal/repository"
	"codewithumam-kasir-api/internal/utils"
)

// TODO: optional try to implement partial update (PATCH)
// TODO: implement proper optimistic locking through etag
type ProductService interface {
	FetchProducts() ([]models.Product, error)
	FetchProductByID(id string) (models.Product, error)
	CreateProduct(product models.CreateProductRequest) (models.Product, error)
	UpdateProductByID(id string, product models.UpdateProductRequest) (models.Product, error)
	DeleteProductByID(id string) error
}

type productService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &productService{
		repository: repository,
	}
}

func (s *productService) FetchProducts() ([]models.Product, error) {
	entities, err := s.repository.FindProducts()
	if err != nil {
		return nil, err
	}

	products := []models.Product{}
	for _, entity := range entities {
		products = append(products, *entity.ToModel())
	}
	return products, nil
}

func (s *productService) FetchProductByID(id string) (models.Product, error) {
	entity, err := s.repository.FindProductByID(utils.DecodeBase62(id))
	if err != nil {
		return models.Product{}, err
	}
	return *entity.ToModel(), nil
}

func (s *productService) CreateProduct(request models.CreateProductRequest) (models.Product, error) {
	entity, err := s.repository.InsertProduct(*request.ToEntity())
	if err != nil {
		return models.Product{}, err
	}
	return *entity.ToModel(), nil
}

func (s *productService) UpdateProductByID(id string, request models.UpdateProductRequest) (models.Product, error) {
	entity, err := s.repository.UpdateProductByID(utils.DecodeBase62(id), *request.ToEntity())
	if err != nil {
		return models.Product{}, err
	}
	return *entity.ToModel(), nil
}

func (s *productService) DeleteProductByID(id string) error {
	return s.repository.DeleteProductByID(utils.DecodeBase62(id))
}
