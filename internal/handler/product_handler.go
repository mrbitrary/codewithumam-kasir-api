package handler

import (
	"codewithumam-kasir-api/internal/models"
	"codewithumam-kasir-api/internal/service"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GET /api/products
func (h *ProductHandler) FetchProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.FetchProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to fetch products"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponseWithItems(products))
	w.WriteHeader(http.StatusOK)
}

// TODO: handle properly if invalid request with correct HTTPStatus
// GET /api/products/{id}
func (h *ProductHandler) FetchProductByID(w http.ResponseWriter, r *http.Request) {
	product, err := h.productService.FetchProductByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to fetch product"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponse(product))
	w.WriteHeader(http.StatusOK)
}

// TODO: handle properly if invalid request with correct HTTPStatus
// POST /api/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var request models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	product, err := h.productService.CreateProduct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to create product"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponse(product))
	w.WriteHeader(http.StatusCreated)
}

// TODO: handle properly if invalid request with correct HTTPStatus
// PUT /api/products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var request models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	product, err := h.productService.UpdateProductByID(r.PathValue("id"), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to update product"))
		return
	}
	_ = json.NewEncoder(w).Encode(models.NewAPIResponse(product))
	w.WriteHeader(http.StatusOK)
}

// TODO: handle properly if invalid request with correct HTTPStatus
// DELETE /api/products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	err := h.productService.DeleteProductByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.NewAPIError(http.StatusInternalServerError, "Failed to delete product"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
