package repository

import (
	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/repository"

	"fmt"
	"github.com/google/uuid"
)

type ProductRepositoryInMemoryImpl struct {
	products []model.ProductEntity
}

func NewProductRepository() repository.ProductRepository {
	return &ProductRepositoryInMemoryImpl{
		products: []model.ProductEntity{},
	}
}

func (r *ProductRepositoryInMemoryImpl) FindProducts() ([]model.ProductEntity, error) {
	return r.products, nil
}

func (r *ProductRepositoryInMemoryImpl) FindProductByID(id string) (model.ProductEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return model.ProductEntity{}, fmt.Errorf("product not found")
	}
	for _, p := range r.products {
		if p.ID == parsedID {
			return p, nil
		}
	}
	return model.ProductEntity{}, fmt.Errorf("product not found")
}

func (r *ProductRepositoryInMemoryImpl) InsertProduct(product model.ProductEntity) (model.ProductEntity, error) {
	r.products = append(r.products, product)
	return product, nil
}

func (r *ProductRepositoryInMemoryImpl) UpdateProductByID(id string, product model.ProductEntity) (model.ProductEntity, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return model.ProductEntity{}, fmt.Errorf("product not found")
	}
	for i, p := range r.products {
		if p.ID == parsedID {
			product.ID = parsedID
			r.products[i] = product
			return product, nil
		}
	}
	return model.ProductEntity{}, fmt.Errorf("product not found")
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
