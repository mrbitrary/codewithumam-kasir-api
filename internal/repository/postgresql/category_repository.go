package repository

import (
	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepositoryPostgreSQLImpl struct {
	connPool *pgxpool.Pool
}

func NewCategoryRepository(connPool *pgxpool.Pool) repository.CategoryRepository {
	return &CategoryRepositoryPostgreSQLImpl{
		connPool: connPool,
	}
}

func (r *CategoryRepositoryPostgreSQLImpl) FindCategories() ([]model.CategoryEntity, error) {
	var categories []model.CategoryEntity
	rows, err := r.connPool.Query(context.Background(), "SELECT id, name, description, created_at, updated_at, deleted_at FROM core.category WHERE deleted_at IS NULL")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.CategoryEntity
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt); err != nil {
			fmt.Println(err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) FindCategoryByID(id string) (model.CategoryEntity, error) {
	var category model.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "SELECT id, name, description, created_at, updated_at, deleted_at, version FROM core.category WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt, &category.Version)
	if err != nil {
		fmt.Println(err)
		return model.CategoryEntity{}, err
	}
	return category, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) FindCategoryByName(name string) (model.CategoryEntity, error) {
	var category model.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "SELECT id, name, description, created_at, updated_at, deleted_at, version FROM core.category WHERE name = $1 AND deleted_at IS NULL", name).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt, &category.Version)
	if err != nil {
		fmt.Println(err)
		return model.CategoryEntity{}, err
	}
	return category, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) InsertCategory(category model.CategoryEntity) (model.CategoryEntity, error) {
	var insertedCategory model.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "INSERT INTO core.category (id, name, description, created_by, updated_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, description, created_by, updated_by, created_at, updated_at, deleted_at, version", category.ID, category.Name, category.Description, category.CreatedBy, category.UpdatedBy).Scan(&insertedCategory.ID, &insertedCategory.Name, &insertedCategory.Description, &insertedCategory.CreatedBy, &insertedCategory.UpdatedBy, &insertedCategory.CreatedAt, &insertedCategory.UpdatedAt, &insertedCategory.DeletedAt, &insertedCategory.Version)
	if err != nil {
		fmt.Println(err)
		return model.CategoryEntity{}, err
	}
	return insertedCategory, nil
}

func (r *CategoryRepositoryPostgreSQLImpl) UpdateCategoryByID(id string, category model.CategoryEntity) (model.CategoryEntity, error) {
	var updatedCategory model.CategoryEntity
	err := r.connPool.QueryRow(context.Background(), "UPDATE core.category SET name = $1, description = $2, updated_by = $3	 WHERE id = $4 AND version = $5 RETURNING id, name, description, updated_by, created_at, updated_at, deleted_at, version", category.Name, category.Description, category.UpdatedBy, id, category.Version).Scan(&updatedCategory.ID, &updatedCategory.Name, &updatedCategory.Description, &updatedCategory.UpdatedBy, &updatedCategory.CreatedAt, &updatedCategory.UpdatedAt, &updatedCategory.DeletedAt, &updatedCategory.Version)
	if err != nil {
		fmt.Println(err)
		return model.CategoryEntity{}, err
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
