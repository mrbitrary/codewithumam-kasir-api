package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"kasir-api/internal/handler"
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
	"kasir-api/internal/service"
)

// TODO: try to implement gzip compressions
// TODO: try to implement rate limit
// TODO: try to implement logging using uber/zap
// TODO: try to implement validation using go-playground/validator
// TODO: try to implement graceful shutdown
// TODO: try to implement proper /heallthprobe, /readinessprobe, /livenessprobe
// TODO: try to integrate with FOSSA, SonarCloud, Snyk, CodeRabbit, CodeCov, GitHub Actions
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.NewAPIResponse(map[string]any{
			"status":    "OK",
			"timestamp": time.Now().Format(time.RFC3339),
		}))
	})

	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	mux.HandleFunc("GET /api/categories", categoryHandler.FetchCategories)
	mux.HandleFunc("GET /api/categories/{id}", categoryHandler.FetchCategoryByID)
	mux.HandleFunc("POST /api/categories", categoryHandler.CreateCategory)
	mux.HandleFunc("PUT /api/categories/{id}", categoryHandler.UpdateCategory)
	mux.HandleFunc("DELETE /api/categories/{id}", categoryHandler.DeleteCategory)

	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)
	mux.HandleFunc("GET /api/products", productHandler.FetchProducts)
	mux.HandleFunc("GET /api/products/{id}", productHandler.FetchProductByID)
	mux.HandleFunc("POST /api/products", productHandler.CreateProduct)
	mux.HandleFunc("PUT /api/products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /api/products/{id}", productHandler.DeleteProduct)

	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(err)
	}
}
