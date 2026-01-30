package repository

import (
	"codewithumam-kasir-api/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepositoryPostgreSQLImpl struct {
	connPool *pgxpool.Pool
}

func NewCategoryRepositoryPostgreSQLImpl(connPool *pgxpool.Pool) CategoryRepository {
	return &CategoryRepositoryPostgreSQLImpl{
		connPool: connPool,
	}
}

func (r *CategoryRepositoryPostgreSQLImpl) FindCategories() ([]models.CategoryEntity, error) {
	var categories []models.CategoryEntity
	rows, err := r.connPool.Query(context.Background(), "SELECT id, name, description FROM core.category WHERE deleted_at IS NULL")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.CategoryEntity
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			fmt.Println(err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) FindCategoryByID(id string) (models.CategoryEntity, error) {
	var category models.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "SELECT id, name, description FROM core.category WHERE id = $1 AND deleted_at IS NULL", id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		fmt.Println(err)
		return models.CategoryEntity{}, err
	}
	return category, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) FindCategoryByName(name string) (models.CategoryEntity, error) {
	var category models.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "SELECT id, name, description FROM core.category WHERE name = $1 AND deleted_at IS NULL", name).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		fmt.Println(err)
		return models.CategoryEntity{}, err
	}
	return category, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) InsertCategory(category models.CategoryEntity) (models.CategoryEntity, error) {
	var insertedCategory models.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "INSERT INTO core.category (id, name, description, created_by, updated_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, created_by, updated_by", category.ID, category.Name, category.Description, category.CreatedBy, category.UpdatedBy).Scan(&insertedCategory.ID, &insertedCategory.Name, &insertedCategory.Description, &insertedCategory.CreatedBy, &insertedCategory.UpdatedBy)
	if err != nil {
		fmt.Println(err)
		return models.CategoryEntity{}, err
	}
	return insertedCategory, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) UpdateCategoryByID(id string, category models.CategoryEntity) (models.CategoryEntity, error) {
	var updatedCategory models.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "UPDATE core.category SET name = $1, description = $2, updated_by = $3	 WHERE id = $4 RETURNING id, name, description, updated_by", category.Name, category.Description, category.UpdatedBy, id).Scan(&updatedCategory.ID, &updatedCategory.Name, &updatedCategory.Description, &updatedCategory.UpdatedBy)
	if err != nil {
		fmt.Println(err)
		return models.CategoryEntity{}, err
	}
	return updatedCategory, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) DeleteCategoryByID(id string) error {
	_, err := r.connPool.Exec(context.Background(), "UPDATE core.category SET deleted_at = NOW(), updated_at = NOW(), updated_by = $1 WHERE id = $2", "USER", id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}