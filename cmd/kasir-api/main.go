package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strings"
	"os"
	"context"

	"github.com/spf13/viper"	
	"github.com/jackc/pgx/v5/pgxpool"

	"codewithumam-kasir-api/config"
	"codewithumam-kasir-api/internal/handler"
	"codewithumam-kasir-api/internal/models"
	"codewithumam-kasir-api/internal/repository"
	"codewithumam-kasir-api/internal/service"
)

// TODO: try to implement gzip compressions
// TODO: try to implement rate limit
// TODO: try to implement logging using uber/zap
// TODO: try to implement validation using go-playground/validator
// TODO: try to implement graceful shutdown
// TODO: try to implement proper /heallthprobe, /readinessprobe, /livenessprobe
// TODO: try to integrate with FOSSA, SonarCloud, Snyk, CodeRabbit, CodeCov, GitHub Actions
func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := &config.Config {
		Port: viper.GetString("PORT"),
		Postgres: config.PostgresConfig {
			Host: viper.GetString("POSTGRES_HOST"),
			Port: viper.GetString("POSTGRES_PORT"),
			User: viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName: viper.GetString("POSTGRES_DB_NAME"),
			SSLMode: viper.GetString("POSTGRES_SSL_MODE"),
			MaxConns: viper.GetInt32("POSTGRES_MAX_CONNS"),
			MaxIdleConnTime: viper.GetDuration("POSTGRES_MAX_IDLE_CONN_TIME"),
			PingTimeout: viper.GetDuration("POSTGRES_PING_TIMEOUT"),
		},
	}

	// postgres://username:password@localhost:5432/database_name?sslmode=require
	pgConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",config.Postgres.User, config.Postgres.Password, config.Postgres.Host, config.Postgres.Port, config.Postgres.DBName, config.Postgres.SSLMode)
	pgxConfig, err := pgxpool.ParseConfig(pgConn)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse postgres config: %v\n", err))
	}

	pgxConfig.MaxConns = config.Postgres.MaxConns
	pgxConfig.MaxConnIdleTime = config.Postgres.MaxIdleConnTime
	pgxConfig.PingTimeout = config.Postgres.PingTimeout

	db, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.NewAPIResponse(map[string]any{
			"status":    "OK",
			"timestamp": time.Now().Format(time.RFC3339),
		}))
	})

	categoryRepository := repository.NewCategoryRepositoryPostgreSQLImpl(db)
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
		panic(err)
	}
}
