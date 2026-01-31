package repository

import (
	"codewithumam-kasir-api/internal/models"
	"codewithumam-kasir-api/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepositoryPostgreSQLImpl struct {
	connPool *pgxpool.Pool
}

func NewProductRepository(connPool *pgxpool.Pool) repository.ProductRepository {
	return &ProductRepositoryPostgreSQLImpl{
		connPool: connPool,
	}
}

func (r *ProductRepositoryPostgreSQLImpl) FindProducts() ([]models.ProductEntity, error) {
	var products []models.ProductEntity
	query := `
		SELECT 
			p.id, p.created_at, p.created_by, p.updated_at, p.updated_by,
			p.name, p.stock, p.price_amount, p.category_id,
			COALESCE(c.name, '') as category_name
		FROM core.product p
		LEFT JOIN core.category c ON p.category_id = c.id AND c.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
	`
	rows, err := r.connPool.Query(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.ProductEntity
		if err := rows.Scan(
			&product.ID, &product.CreatedAt, &product.CreatedBy, &product.UpdatedAt, &product.UpdatedBy,
			&product.Name, &product.Stocks, &product.Price, &product.CategoryID,
			&product.CategoryName,
		); err != nil {
			fmt.Println(err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryPostgreSQLImpl) FindProductByID(id string) (models.ProductEntity, error) {
	var product models.ProductEntity
	query := `
		SELECT 
			p.id, p.version, p.created_at, p.created_by, p.updated_at, p.updated_by, p.deleted_at,
			p.name, p.stock, p.price_amount, p.category_id,
			COALESCE(c.name, '') as category_name
		FROM core.product p
		LEFT JOIN core.category c ON p.category_id = c.id AND c.deleted_at IS NULL
		WHERE p.id = $1
	`
	err := r.connPool.QueryRow(context.Background(), query, id).Scan(
		&product.ID, &product.Version, &product.CreatedAt, &product.CreatedBy, &product.UpdatedAt, &product.UpdatedBy, &product.DeletedAt,
		&product.Name, &product.Stocks, &product.Price, &product.CategoryID,
		&product.CategoryName,
	)
	if err != nil {
		fmt.Println(err)
		return models.ProductEntity{}, err
	}
	return product, nil
}

func (r *ProductRepositoryPostgreSQLImpl) InsertProduct(product models.ProductEntity) (models.ProductEntity, error) {
	var insertedProduct models.ProductEntity
	query := `
		WITH category_lookup AS (
			SELECT id FROM core.category WHERE lower(name) = lower($5) AND deleted_at IS NULL
		)
		INSERT INTO core.product (
			id, name, stock, price_amount, price_scale, currency, category_id,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, 0, 'IDR', (SELECT id FROM category_lookup), $6, $7
		)
		RETURNING 
			id, version, created_at, created_by, updated_at, updated_by,
			name, stock, price_amount, category_id
	`
	err := r.connPool.QueryRow(
		context.Background(), 
		query, 
		product.ID, product.Name, product.Stocks, product.Price, product.CategoryName, 
		product.CreatedBy, product.UpdatedBy,
	).Scan(
		&insertedProduct.ID, &insertedProduct.Version, &insertedProduct.CreatedAt, &insertedProduct.CreatedBy, &insertedProduct.UpdatedAt, &insertedProduct.UpdatedBy,
		&insertedProduct.Name, &insertedProduct.Stocks, &insertedProduct.Price, &insertedProduct.CategoryID,
	)
	
	if err != nil {
		fmt.Println(err)
		return models.ProductEntity{}, err
	}
	
	insertedProduct.CategoryName = product.CategoryName 
	
	return insertedProduct, nil
}

func (r *ProductRepositoryPostgreSQLImpl) UpdateProductByID(id string, product models.ProductEntity) (models.ProductEntity, error) {
	var updatedProduct models.ProductEntity
	query := `
		WITH category_lookup AS (
			SELECT id FROM core.category WHERE lower(name) = lower($4) AND deleted_at IS NULL
		)
		UPDATE core.product 
		SET 
			name = $1, 
			stock = $2,
			price_amount = $3,
			category_id = (SELECT id FROM category_lookup),
			updated_by = $5
		WHERE id = $6 AND version = $7 AND deleted_at IS NULL
		RETURNING 
			id, version, created_at, created_by, updated_at, updated_by, deleted_at,
			name, stock, price_amount, category_id
	`
	
	err := r.connPool.QueryRow(
		context.Background(), 
		query,
		product.Name, product.Stocks, product.Price, product.CategoryName,
		product.UpdatedBy, id, product.Version,
	).Scan(
		&updatedProduct.ID, &updatedProduct.Version, &updatedProduct.CreatedAt, &updatedProduct.CreatedBy, &updatedProduct.UpdatedAt, &updatedProduct.UpdatedBy, &updatedProduct.DeletedAt,
		&updatedProduct.Name, &updatedProduct.Stocks, &updatedProduct.Price, &updatedProduct.CategoryID,
	)

	if err != nil {
		fmt.Println(err)
		return models.ProductEntity{}, err
	}
	
	updatedProduct.CategoryName = product.CategoryName
	return updatedProduct, nil
}

func (r *ProductRepositoryPostgreSQLImpl) DeleteProductByID(id string) error {
	_, err := r.connPool.Exec(context.Background(), "UPDATE core.product SET deleted_at = NOW(), updated_at = NOW(), updated_by = $1 WHERE id = $2", "USER", id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}