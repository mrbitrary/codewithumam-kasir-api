package repository

import (
	"context"
	"fmt"
	"time"

	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepositoryPostgreSQLImpl struct {
	connPool *pgxpool.Pool
}

func NewTransactionRepository(connPool *pgxpool.Pool) repository.TransactionRepository {
	return &TransactionRepositoryPostgreSQLImpl{
		connPool: connPool,
	}
}

func (r *TransactionRepositoryPostgreSQLImpl) CreateTransaction(tx model.TransactionEntity, details []model.TransactionDetailEntity) (model.TransactionEntity, error) {
	ctx := context.Background()
	conn, err := r.connPool.Begin(ctx)
	if err != nil {
		return model.TransactionEntity{}, err
	}
	defer func() {
		_ = conn.Rollback(ctx)
	}()

	txQuery := `
		INSERT INTO core.transaction (
			id, total_items, total_price_amount, total_price_scale, currency, 
			created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = conn.Exec(ctx, txQuery,
		tx.ID, tx.TotalItems, tx.TotalPriceAmount, tx.TotalPriceScale, tx.Currency,
		tx.CreatedBy, tx.UpdatedBy,
	)
	if err != nil {
		return model.TransactionEntity{}, fmt.Errorf("failed to insert transaction: %w", err)
	}

	detailQuery := `
		INSERT INTO core.transaction_detail (
			id, transaction_id, product_id, product_name, category_id, category_name,
			price_amount, price_scale, currency,
			quantity, total_price_amount, total_price_scale, 
			created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	stockQuery := `
		UPDATE core.product 
		SET stock = stock - $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	for _, d := range details {
		_, err = conn.Exec(ctx, detailQuery,
			d.ID, d.TransactionID, d.ProductID, d.ProductName, d.CategoryID, d.CategoryName,
			d.PriceAmount, d.PriceScale, d.Currency,
			d.Quantity, d.TotalPriceAmount, d.TotalPriceScale,
			d.CreatedBy, d.UpdatedBy,
		)
		if err != nil {
			return model.TransactionEntity{}, fmt.Errorf("failed to insert transaction detail: %w", err)
		}

		if d.ProductID != nil {
			cmd, err := conn.Exec(ctx, stockQuery, d.Quantity, d.CreatedBy, d.ProductID)
			if err != nil {
				return model.TransactionEntity{}, fmt.Errorf("failed to update stock: %w", err)
			}
			if cmd.RowsAffected() == 0 {
				return model.TransactionEntity{}, fmt.Errorf("failed to update stock: product %s not found or deleted", d.ProductName)
			}
		}
	}

	if err := conn.Commit(ctx); err != nil {
		return model.TransactionEntity{}, err
	}

	return tx, nil
}

func (r *TransactionRepositoryPostgreSQLImpl) GetReportStats(startDate, endDate time.Time) (model.ReportResponse, error) {
	ctx := context.Background()
	var report model.ReportResponse

	query := `
		SELECT 
			COALESCE(SUM(total_revenue), 0), 
			COALESCE(SUM(total_transactions), 0)
		FROM core.transaction_summary_daily
		WHERE report_date >= $1 AND report_date <= $2
	`
	err := r.connPool.QueryRow(ctx, query, startDate, endDate).Scan(&report.TotalRevenue.Amount, &report.TotalTransactions)
	if err != nil {
		return report, err
	}
	report.TotalRevenue.Currency = "IDR" // Default
	report.TotalRevenue.Display = float64(report.TotalRevenue.Amount)

	topItemsQuery := `
		SELECT p.name, SUM(s.total_sold) as total_qty
		FROM core.sales_summary_daily s
		JOIN core.product p ON s.product_id = p.id
		WHERE s.report_date >= $1 AND s.report_date <= $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 5
	`
	rows, err := r.connPool.Query(ctx, topItemsQuery, startDate, endDate)
	if err != nil {
		return report, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.PopularItem
		if err := rows.Scan(&item.Name, &item.TotalSoldQty); err != nil {
			return report, err
		}
		report.TopPopularItems = append(report.TopPopularItems, item)
	}

	topCatsQuery := `
		SELECT COALESCE(c.name, 'Uncategorized'), SUM(s.total_sold) as total_qty
		FROM core.sales_summary_daily s
		LEFT JOIN core.category c ON s.category_id = c.id
		WHERE s.report_date >= $1 AND s.report_date <= $2
		GROUP BY c.name, s.category_id
		ORDER BY total_qty DESC
		LIMIT 5
	`
	rows2, err := r.connPool.Query(ctx, topCatsQuery, startDate, endDate)
	if err != nil {
		return report, err
	}
	defer rows2.Close()

	for rows2.Next() {
		var cat model.PopularCategory
		if err := rows2.Scan(&cat.Name, &cat.TotalSoldQty); err != nil {
			return report, err
		}
		report.TopPopularCategories = append(report.TopPopularCategories, cat)
	}

	return report, nil
}

func (r *TransactionRepositoryPostgreSQLImpl) GetMostPopularCategory(startDate, endDate time.Time) (model.PopularCategory, error) {
	ctx := context.Background()
	var category model.PopularCategory
	query := `
		SELECT COALESCE(c.name, 'Uncategorized'), SUM(s.total_sold) as total_qty
		FROM core.sales_summary_daily s
		LEFT JOIN core.category c ON s.category_id = c.id
		WHERE s.report_date >= $1 AND s.report_date <= $2
		GROUP BY c.name, s.category_id
		ORDER BY total_qty DESC
		LIMIT 1
	`
	err := r.connPool.QueryRow(ctx, query, startDate, endDate).Scan(&category.Name, &category.TotalSoldQty)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *TransactionRepositoryPostgreSQLImpl) GetMostPopularProduct(startDate, endDate time.Time) (model.PopularItem, error) {
	ctx := context.Background()
	var product model.PopularItem
	query := `
		SELECT p.name, SUM(s.total_sold) as total_qty
		FROM core.sales_summary_daily s
		JOIN core.product p ON s.product_id = p.id
		WHERE s.report_date >= $1 AND s.report_date <= $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`
	err := r.connPool.QueryRow(ctx, query, startDate, endDate).Scan(&product.Name, &product.TotalSoldQty)
	if err != nil {
		return product, err
	}
	return product, nil
}
