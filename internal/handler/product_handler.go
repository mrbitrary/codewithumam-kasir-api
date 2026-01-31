package handler

import (
	"codewithumam-kasir-api/internal/model"
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
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusInternalServerError, "Failed to fetch products"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.NewAPIResponseWithItems(products))
}

// TODO: handle properly if invalid request with correct HTTPStatus
// GET /api/products/{id}
func (h *ProductHandler) FetchProductByID(w http.ResponseWriter, r *http.Request) {
	product, err := h.productService.FetchProductByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusInternalServerError, "Failed to fetch product"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.NewAPIResponse(product))
}

// TODO: handle properly if invalid request with correct HTTPStatus
// POST /api/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var request model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	product, err := h.productService.CreateProduct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusInternalServerError, "Failed to create product"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(model.NewAPIResponse(product))
}

// TODO: handle properly if invalid request with correct HTTPStatus
// PUT /api/products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var request model.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusBadRequest, "Invalid request body"))
		return
	}
	product, err := h.productService.UpdateProductByID(r.PathValue("id"), request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusInternalServerError, "Failed to update product"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.NewAPIResponse(product))
}

// TODO: handle properly if invalid request with correct HTTPStatus
// DELETE /api/products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	err := h.productService.DeleteProductByID(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.NewAPIError(http.StatusInternalServerError, "Failed to delete product"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
