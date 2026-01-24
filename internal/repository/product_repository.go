package repository

import (
	"codewithumam-kasir-api/internal/models"
	"fmt"
	"github.com/google/uuid"
)

type ProductRepository interface {
	FindProducts() ([]models.ProductEntity, error)
	FindProductByID(id string) (models.ProductEntity, error)
	InsertProduct(product models.ProductEntity) (models.ProductEntity, error)
	UpdateProductByID(id string, product models.ProductEntity) (models.ProductEntity, error)
	DeleteProductByID(id string) error
}

type ProductRepositoryInMemoryImpl struct {
	products []models.ProductEntity
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryInMemoryImpl{
		products: []models.ProductEntity{},
	}
}

func (r *ProductRepositoryInMemoryImpl) FindProducts() ([]models.ProductEntity, error) {
	return r.products, nil
}

func (r *ProductRepositoryInMemoryImpl) FindProductByID(id string) (models.ProductEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.ProductEntity{}, fmt.Errorf("product not found")
	}
	for _, p := range r.products {
		if p.ID == parsedID {
			return p, nil
		}
	}
	return models.ProductEntity{}, fmt.Errorf("product not found")
}

func (r *ProductRepositoryInMemoryImpl) InsertProduct(product models.ProductEntity) (models.ProductEntity, error) {
	r.products = append(r.products, product)
	return product, nil
}

func (r *ProductRepositoryInMemoryImpl) UpdateProductByID(id string, product models.ProductEntity) (models.ProductEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.ProductEntity{}, fmt.Errorf("product not found")
	}
	for i, p := range r.products {
		if p.ID == parsedID {
			product.ID = parsedID
			r.products[i] = product
			return product, nil
		}
	}
	return models.ProductEntity{}, fmt.Errorf("product not found")
}

func (r *ProductRepositoryInMemoryImpl) DeleteProductByID(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("product not found")
	}
	for i, p := range r.products {
		if p.ID == parsedID {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product not found")
}
