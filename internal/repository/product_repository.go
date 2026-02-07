package repository

import (
	"codewithumam-kasir-api/internal/model"
)

type ProductRepository interface {
	FindProducts() ([]model.ProductEntity, error)
	FindProductByID(id string) (model.ProductEntity, error)
	FindProductsByNameAndActiveStatus(name string, activeStatus *bool) ([]model.ProductEntity, error)
	InsertProduct(product model.ProductEntity) (model.ProductEntity, error)
	UpdateProductByID(id string, product model.ProductEntity) (model.ProductEntity, error)
	DeleteProductByID(id string) error
}
