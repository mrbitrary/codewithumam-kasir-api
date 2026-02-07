package service

import (
	"errors"
	"time"

	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/repository"
	"codewithumam-kasir-api/internal/utils"
	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(req model.CreateTransactionRequest) (model.Transaction, error)
	FetchReport(startDateStr, endDateStr, period string) (model.ReportResponse, error)
	FetchMostPopularCategory(startDateStr, endDateStr string) (model.PopularCategory, error)
	FetchMostPopularProduct(startDateStr, endDateStr string) (model.PopularItem, error)
}

type TransactionServiceImpl struct {
	txRepo      repository.TransactionRepository
	productRepo repository.ProductRepository
}

func NewTransactionService(txRepo repository.TransactionRepository, productRepo repository.ProductRepository) TransactionService {
	return &TransactionServiceImpl{
		txRepo:      txRepo,
		productRepo: productRepo,
	}
}

func (s *TransactionServiceImpl) CreateTransaction(req model.CreateTransactionRequest) (model.Transaction, error) {
	if len(req.Items) == 0 {
		return model.Transaction{}, errors.New("transaction must have at least one item")
	}

	txID, _ := uuid.NewV7()
	var totalItems int
	var totalPriceAmount int64
	var details []model.TransactionDetailEntity

	currency := "IDR"
	scale := 0

	for _, item := range req.Items {
		product, err := s.productRepo.FindProductByID(utils.DecodeBase62(item.ProductID))
		if err != nil {
			return model.Transaction{}, err
		}

		if product.Stocks < item.Quantity {
			return model.Transaction{}, errors.New("insufficient stock for product: " + product.Name)
		}

		detailID, _ := uuid.NewV7()
		itemTotalPrice := product.Price * int64(item.Quantity)

		detail := model.TransactionDetailEntity{
			ID:                detailID,
			TransactionID:     txID,
			ProductID:         &product.ID,
			ProductName:       product.Name,
			CategoryID:        product.CategoryID,
			CategoryName:      product.CategoryName,
			PriceAmount:       product.Price,
			PriceScale:        scale,
			PriceDisplay:      float64(product.Price), // Simple conversion for now
			Currency:          currency,
			Quantity:          item.Quantity,
			TotalPriceAmount:  itemTotalPrice,
			TotalPriceScale:   scale,
			TotalPriceDisplay: float64(itemTotalPrice),
			CreatedBy:         "USER",
			UpdatedBy:         "USER",
		}

		details = append(details, detail)
		totalItems += item.Quantity
		totalPriceAmount += itemTotalPrice
	}

	txEntity := model.TransactionEntity{
		ID:                txID,
		TotalItems:        totalItems,
		TotalPriceAmount:  totalPriceAmount,
		TotalPriceScale:   scale,
		TotalPriceDisplay: float64(totalPriceAmount),
		Currency:          currency,
		CreatedBy:         "USER",
		UpdatedBy:         "USER",
		CreatedAt:         time.Now(),
	}

	createdTx, err := s.txRepo.CreateTransaction(txEntity, details)
	if err != nil {
		return model.Transaction{}, err
	}

	result := createdTx.ToModel()
	// Populate details for the response
	for _, d := range details {
		result.Details = append(result.Details, *d.ToModel())
	}

	return *result, nil
}

func (s *TransactionServiceImpl) FetchReport(startDateStr, endDateStr, period string) (model.ReportResponse, error) {
	startDate, endDate := s.parseDateRange(startDateStr, endDateStr, period)
	if startDate.After(endDate) {
		return model.ReportResponse{}, errors.New("startDate cannot be after endDate")
	}
	return s.txRepo.GetReportStats(startDate, endDate)
}

func (s *TransactionServiceImpl) FetchMostPopularCategory(startDateStr, endDateStr string) (model.PopularCategory, error) {
	startDate, endDate := s.parseDateRange(startDateStr, endDateStr, "")
	if startDate.After(endDate) {
		return model.PopularCategory{}, errors.New("startDate cannot be after endDate")
	}
	return s.txRepo.GetMostPopularCategory(startDate, endDate)
}

func (s *TransactionServiceImpl) FetchMostPopularProduct(startDateStr, endDateStr string) (model.PopularItem, error) {
	startDate, endDate := s.parseDateRange(startDateStr, endDateStr, "")
	if startDate.After(endDate) {
		return model.PopularItem{}, errors.New("startDate cannot be after endDate")
	}
	return s.txRepo.GetMostPopularProduct(startDate, endDate)
}

func (s *TransactionServiceImpl) parseDateRange(startDateStr, endDateStr, period string) (time.Time, time.Time) {
	now := time.Now()
	// Set to start of day and end of day
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	if period != "" {
		switch period {
		case "today":
			return todayStart, todayEnd
		case "yesterday":
			y := todayStart.AddDate(0, 0, -1)
			return y, time.Date(y.Year(), y.Month(), y.Day(), 23, 59, 59, 999999999, y.Location())
		case "last-week":
			// Monday-Sunday of previous week
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			mondayThisWeek := todayStart.AddDate(0, 0, -(weekday - 1))
			lastSunday := mondayThisWeek.AddDate(0, 0, -1)
			lastSundayEnd := time.Date(lastSunday.Year(), lastSunday.Month(), lastSunday.Day(), 23, 59, 59, 999999999, lastSunday.Location())
			lastMonday := lastSunday.AddDate(0, 0, -6)
			return lastMonday, lastSundayEnd
		case "last-month":
			// 1st to end-of-month of previous month
			firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			lastMonthEnd := firstOfThisMonth.Add(-time.Nanosecond)
			lastMonthStart := time.Date(lastMonthEnd.Year(), lastMonthEnd.Month(), 1, 0, 0, 0, 0, now.Location())
			return lastMonthStart, lastMonthEnd
		case "week-to-date":
			// Assuming week starts on Monday
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			return todayStart.AddDate(0, 0, -(weekday - 1)), todayEnd
		case "month-to-date":
			return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), todayEnd
		case "year-to-date":
			return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()), todayEnd
		}
	}

	start, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		start = todayStart
	}
	end, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		end = todayEnd
	} else {
		// Ensure end is end of day
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	}

	return start, end
}
