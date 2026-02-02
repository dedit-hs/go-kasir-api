package main

import (
	"encoding/json"
	"fmt"
	"go-kasir-api/database"
	"go-kasir-api/handlers"
	"go-kasir-api/repositories"
	"go-kasir-api/services"
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func maskConnectionString(conn string) string {
	if len(conn) > 20 {
		return conn[:20] + "..." + conn[len(conn)-10:]
	}
	return "***"
}

func main() {
	viper.AutomaticEnv()

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	log.Printf("Configuration loaded - Port: %s, DB_CONN: %s", config.Port, maskConnectionString(config.DBConn))

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	defer db.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	productRepository := repositories.NewProductRepositories(db)
	productService := services.NewProductServices(productRepository)
	productHandler := handlers.NewProductHandlers(productService)

	categoryRepository := repositories.NewCategoryRepositories(db)
	categoryService := services.NewCategoryServices(categoryRepository)
	categoryHandler := handlers.NewCategoryHandlers(categoryService)

	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/{id}", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/{id}", categoryHandler.HandleCategoryByID)

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server started on :" + addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
