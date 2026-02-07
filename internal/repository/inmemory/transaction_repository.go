package repository

import (
	"time"

	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/repository"
)

type TransactionRepositoryInMemoryImpl struct {
	transactions []model.TransactionEntity
	details      []model.TransactionDetailEntity
	productRepo  repository.ProductRepository
}

func NewTransactionRepository(productRepo repository.ProductRepository) repository.TransactionRepository {
	return &TransactionRepositoryInMemoryImpl{
		transactions: []model.TransactionEntity{},
		details:      []model.TransactionDetailEntity{},
		productRepo:  productRepo,
	}
}

func (r *TransactionRepositoryInMemoryImpl) CreateTransaction(tx model.TransactionEntity, details []model.TransactionDetailEntity) (model.TransactionEntity, error) {
	for _, d := range details {
		if d.ProductID != nil {
			p, _ := r.productRepo.FindProductByID(d.ProductID.String())
			p.Stocks -= d.Quantity
			p.UpdatedAt = time.Now()
			_, _ = r.productRepo.UpdateProductByID(p.ID.String(), p)
		}
	}

	r.transactions = append(r.transactions, tx)
	r.details = append(r.details, details...)

	return tx, nil
}

func (r *TransactionRepositoryInMemoryImpl) GetReportStats(startDate, endDate time.Time) (model.ReportResponse, error) {
	var totalRevenue int64
	var totalTransactions int
	for _, tx := range r.transactions {
		if (tx.CreatedAt.After(startDate) || tx.CreatedAt.Equal(startDate)) &&
			(tx.CreatedAt.Before(endDate) || tx.CreatedAt.Equal(endDate)) {
			totalRevenue += tx.TotalPriceAmount
			totalTransactions++
		}
	}
	return model.ReportResponse{
		TotalRevenue: model.Price{
			Amount:   totalRevenue,
			Currency: "IDR",
			Display:  float64(totalRevenue),
		},
		TotalTransactions: totalTransactions,
	}, nil
}

func (r *TransactionRepositoryInMemoryImpl) GetMostPopularCategory(startDate, endDate time.Time) (model.PopularCategory, error) {
	return model.PopularCategory{}, nil
}

func (r *TransactionRepositoryInMemoryImpl) GetMostPopularProduct(startDate, endDate time.Time) (model.PopularItem, error) {
	return model.PopularItem{}, nil
}
