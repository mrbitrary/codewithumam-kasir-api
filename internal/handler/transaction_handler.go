package handler

import (
	"encoding/json"
	"net/http"

	"codewithumam-kasir-api/internal/model"
	"codewithumam-kasir-api/internal/service"
)

type TransactionHandler struct {
	txService service.TransactionService
}

func NewTransactionHandler(txService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		txService: txService,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req model.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	tx, err := h.txService.CreateTransaction(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusBadRequest, err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(model.NewAPIResponse(tx))
}

func (h *TransactionHandler) FetchReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	period := r.URL.Query().Get("period")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	path := r.URL.Path
	if period == "" {
		switch path {
		case "/api/reports/today":
			period = "today"
		case "/api/reports/yesterday":
			period = "yesterday"
		case "/api/reports/last-week":
			period = "last-week"
		case "/api/reports/last-month":
			period = "last-month"
		case "/api/reports/week-to-date":
			period = "week-to-date"
		case "/api/reports/month-to-date":
			period = "month-to-date"
		case "/api/reports/year-to-date":
			period = "year-to-date"
		}
	}

	report, err := h.txService.FetchReport(startDate, endDate, period)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "startDate cannot be after endDate" {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(status, err.Error()))
		return
	}

	_ = json.NewEncoder(w).Encode(model.NewAPIResponse(report))
}

func (h *TransactionHandler) FetchPopularCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	category, err := h.txService.FetchMostPopularCategory(startDate, endDate)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "startDate cannot be after endDate" {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(status, err.Error()))
		return
	}

	_ = json.NewEncoder(w).Encode(model.NewAPIResponse(category))
}

func (h *TransactionHandler) FetchPopularProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	product, err := h.txService.FetchMostPopularProduct(startDate, endDate)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "startDate cannot be after endDate" {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(status, err.Error()))
		return
	}

	_ = json.NewEncoder(w).Encode(model.NewAPIResponse(product))
}
