package repository

import (
	"codewithumam-kasir-api/internal/models"
)

type ProductRepository interface {
	FindProducts() ([]models.ProductEntity, error)
	FindProductByID(id string) (models.ProductEntity, error)
	InsertProduct(product models.ProductEntity) (models.ProductEntity, error)
	UpdateProductByID(id string, product models.ProductEntity) (models.ProductEntity, error)
	DeleteProductByID(id string) error
}
