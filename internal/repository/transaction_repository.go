package repository

import (
	"time"

	"codewithumam-kasir-api/internal/model"
)

type TransactionRepository interface {
	CreateTransaction(tx model.TransactionEntity, details []model.TransactionDetailEntity) (model.TransactionEntity, error)
	GetReportStats(startDate, endDate time.Time) (model.ReportResponse, error)
	GetMostPopularCategory(startDate, endDate time.Time) (model.PopularCategory, error)
	GetMostPopularProduct(startDate, endDate time.Time) (model.PopularItem, error)
}
