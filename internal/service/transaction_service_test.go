package service

import (
	"codewithumam-kasir-api/internal/mock"
	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	mockTxRepo := new(mock.MockTransactionRepository)
	mockProductRepo := new(mock.MockProductRepository)
	service := NewTransactionService(mockTxRepo, mockProductRepo)

	productID, _ := uuid.NewV7()
	product := model.ProductEntity{
		ID:           productID,
		Name:         "Test Product",
		Price:        1000,
		Stocks:       10,
		CategoryName: "Category A",
	}

	req := model.CreateTransactionRequest{
		Items: []model.CreateTransactionItemRequest{
			{
				ProductID: utils.EncodeBase62(productID.String()),
				Quantity:  2,
			},
		},
	}

	mockProductRepo.On("FindProductByID", productID.String()).Return(product, nil)
	mockTxRepo.On("CreateTransaction", testifyMock.Anything, testifyMock.Anything).Return(model.TransactionEntity{ID: productID}, nil)

	tx, err := service.CreateTransaction(req)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	mockProductRepo.AssertExpectations(t)
	mockTxRepo.AssertExpectations(t)
}

func TestTransactionService_CreateTransaction_InsufficientStock(t *testing.T) {
	mockTxRepo := new(mock.MockTransactionRepository)
	mockProductRepo := new(mock.MockProductRepository)
	service := NewTransactionService(mockTxRepo, mockProductRepo)

	productID, _ := uuid.NewV7()
	product := model.ProductEntity{
		ID:     productID,
		Name:   "Test Product",
		Stocks: 1,
	}

	req := model.CreateTransactionRequest{
		Items: []model.CreateTransactionItemRequest{
			{
				ProductID: utils.EncodeBase62(productID.String()),
				Quantity:  2,
			},
		},
	}

	mockProductRepo.On("FindProductByID", productID.String()).Return(product, nil)

	_, err := service.CreateTransaction(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
}

func TestTransactionService_FetchReport(t *testing.T) {
	mockTxRepo := new(mock.MockTransactionRepository)
	mockProductRepo := new(mock.MockProductRepository)
	service := NewTransactionService(mockTxRepo, mockProductRepo)

	mockTxRepo.On("GetReportStats", testifyMock.Anything, testifyMock.Anything).Return(model.ReportResponse{TotalTransactions: 5}, nil)

	resp, err := service.FetchReport("", "", "today")

	assert.NoError(t, err)
	assert.Equal(t, 5, resp.TotalTransactions)
}

func TestTransactionService_Reports(t *testing.T) {
	mockTxRepo := new(mock.MockTransactionRepository)
	mockProductRepo := new(mock.MockProductRepository)
	service := NewTransactionService(mockTxRepo, mockProductRepo)

	mockTxRepo.On("GetMostPopularCategory", testifyMock.Anything, testifyMock.Anything).Return(model.PopularCategory{Name: "Cat"}, nil)
	mockTxRepo.On("GetMostPopularProduct", testifyMock.Anything, testifyMock.Anything).Return(model.PopularItem{Name: "Prod"}, nil)

	cat, _ := service.FetchMostPopularCategory("", "")
	prod, _ := service.FetchMostPopularProduct("", "")

	assert.Equal(t, "Cat", cat.Name)
	assert.Equal(t, "Prod", prod.Name)
}

func TestTransactionService_ParseDateRange(t *testing.T) {
	mockTxRepo := new(mock.MockTransactionRepository)
	mockProductRepo := new(mock.MockProductRepository)
	s := &TransactionServiceImpl{txRepo: mockTxRepo, productRepo: mockProductRepo}

	// Test periods
	periods := []struct {
		period string
		check  func(t *testing.T, start, end time.Time)
	}{
		{
			period: "today",
			check: func(t *testing.T, start, end time.Time) {
				now := time.Now()
				assert.Equal(t, now.Day(), start.Day())
				assert.Equal(t, 0, start.Hour())
				assert.Equal(t, 23, end.Hour())
			},
		},
		{
			period: "yesterday",
			check: func(t *testing.T, start, end time.Time) {
				expected := time.Now().AddDate(0, 0, -1)
				assert.Equal(t, expected.Day(), start.Day())
			},
		},
		{
			period: "last-week",
			check: func(t *testing.T, start, end time.Time) {
				assert.Equal(t, time.Sunday, end.Weekday())
				assert.Equal(t, time.Monday, start.Weekday())
				assert.True(t, end.After(start))
			},
		},
		{
			period: "last-month",
			check: func(t *testing.T, start, end time.Time) {
				assert.Equal(t, 1, start.Day())
				assert.True(t, end.Before(time.Now()))
			},
		},
		{
			period: "week-to-date",
			check: func(t *testing.T, start, end time.Time) {
				assert.Equal(t, time.Monday, start.Weekday())
				assert.Equal(t, time.Now().Day(), end.Day())
			},
		},
		{
			period: "month-to-date",
			check: func(t *testing.T, start, end time.Time) {
				assert.Equal(t, 1, start.Day())
			},
		},
		{
			period: "year-to-date",
			check: func(t *testing.T, start, end time.Time) {
				assert.Equal(t, 1, start.Day())
				assert.Equal(t, 1, int(start.Month()))
			},
		},
	}

	for _, tc := range periods {
		t.Run(tc.period, func(t *testing.T) {
			start, end := s.parseDateRange("", "", tc.period)
			assert.False(t, start.IsZero())
			assert.False(t, end.IsZero())
			tc.check(t, start, end)
		})
	}

	// Test explicit dates
	start, end := s.parseDateRange("2024-01-01", "2024-01-02", "")
	assert.Equal(t, 2024, start.Year())
	assert.Equal(t, 1, int(start.Month()))
	assert.Equal(t, 1, start.Day())
	assert.Equal(t, 2, end.Day())
	assert.Equal(t, 23, end.Hour())
}
func TestTransactionService_FetchReport_InvalidDateRange(t *testing.T) {
	mockTxRepo := new(mock.MockTransactionRepository)
	mockProductRepo := new(mock.MockProductRepository)
	service := NewTransactionService(mockTxRepo, mockProductRepo)

	_, err := service.FetchReport("2024-01-02", "2024-01-01", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "startDate cannot be after endDate")
}
