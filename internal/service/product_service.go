package service

import (
	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/repository"
	"codewithumam-kasir-api/internal/utils"
)

// TODO: optional try to implement partial update (PATCH)
// TODO: implement proper optimistic locking through etag or version checks
// TODO: handle etag properly
type ProductService interface {
	FetchProducts() ([]model.Product, error)
	FetchProductsByNameAndActiveStatus(name string, activeStatus *bool) ([]model.Product, error)
	FetchProductByID(id string) (model.Product, error)
	CreateProduct(product model.CreateProductRequest) (model.Product, error)
	UpdateProductByID(id string, product model.UpdateProductRequest) (model.Product, error)
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

func (s *productService) FetchProducts() ([]model.Product, error) {
	entities, err := s.repository.FindProducts()
	if err != nil {
		return nil, err
	}

	products := []model.Product{}
	for _, entity := range entities {
		products = append(products, *entity.ToModel())
	}
	return products, nil
}

func (s *productService) FetchProductsByNameAndActiveStatus(name string, activeStatus *bool) ([]model.Product, error) {
	entities, err := s.repository.FindProductsByNameAndActiveStatus(name, activeStatus)
	if err != nil {
		return nil, err
	}

	products := []model.Product{}
	for _, entity := range entities {
		products = append(products, *entity.ToModel())
	}
	return products, nil
}

func (s *productService) FetchProductByID(id string) (model.Product, error) {
	entity, err := s.repository.FindProductByID(utils.DecodeBase62(id))
	if err != nil {
		return model.Product{}, err
	}
	return *entity.ToModel(), nil
}

func (s *productService) CreateProduct(request model.CreateProductRequest) (model.Product, error) {
	entity, err := s.repository.InsertProduct(*request.ToEntity())
	if err != nil {
		return model.Product{}, err
	}
	return *entity.ToModel(), nil
}

func (s *productService) UpdateProductByID(id string, request model.UpdateProductRequest) (model.Product, error) {
	entity, err := s.repository.UpdateProductByID(utils.DecodeBase62(id), *request.ToEntity())
	if err != nil {
		return model.Product{}, err
	}
	return *entity.ToModel(), nil
}

func (s *productService) DeleteProductByID(id string) error {
	return s.repository.DeleteProductByID(utils.DecodeBase62(id))
}
