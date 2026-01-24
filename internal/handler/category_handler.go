package handler

import (
	"encoding/json"
	"kasir-api/internal/models"
	"kasir-api/internal/service"
	"net/http"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}

}

// GET /api/categories
func (h *CategoryHandler) FetchCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories, err := h.categoryService.FetchCategories()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to fetch categories"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponseWithItems(categories))
	w.WriteHeader(http.StatusOK)

}

// TODO: handle properly if invalid request with correct HTTPStatus
// GET /api/categories/{id}
func (h *CategoryHandler) FetchCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	category, err := h.categoryService.FetchCategoryByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to fetch category"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponse(category))
	w.WriteHeader(http.StatusOK)

}

// TODO: handle properly if invalid request with correct HTTPStatus
// POST /api/categories
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := models.CreateCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}

	category, err := h.categoryService.CreateCategory(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to create category"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponse(category))
	w.WriteHeader(http.StatusCreated)
}

// TODO: handle properly if invalid request with correct HTTPStatus
// PUT /api/categories/{id}
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	request := models.UpdateCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	category, err := h.categoryService.UpdateCategoryByID(r.PathValue("id"), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to update category"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponse(category))
	w.WriteHeader(http.StatusOK)
}

// TODO: handle properly if invalid request with correct HTTPStatus
// DELETE /api/categories/{id}
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := h.categoryService.DeleteCategoryByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to delete category"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
