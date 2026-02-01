package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"PORT"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var Products []Product = []Product{
	{ID: 1, Name: "Indomie Goreng", Price: 3500, Stock: 40},
	{ID: 2, Name: "Aqua 600ml", Price: 3000, Stock: 24},
	{ID: 3, Name: "Bodrex", Price: 1000, Stock: 20},
	{ID: 4, Name: "Gula 1kg", Price: 15000, Stock: 10},
	{ID: 5, Name: "Pensil", Price: 2000, Stock: 12},
}

var Categories []Category = []Category{
	{ID: 1, Name: "Makanan", Description: "Makanan Ringan"},
	{ID: 2, Name: "Minuman", Description: "Minuman Menyehatkan"},
	{ID: 3, Name: "Obat", Description: "Obat-obatan"},
	{ID: 4, Name: "Sembako", Description: "Sembilan Bahan Pokok"},
	{ID: 5, Name: "Alat Tulis", Description: "Alat Tulis Sekolah"},
}

func UnknownMethod(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
}

func InvalidRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		InvalidRequest(w, r)
		return
	}
	product.ID = len(Products) + 1
	Products = append(Products, product)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func GetProduct(w http.ResponseWriter, r *http.Request, productId int) {
	for _, product := range Products {
		if product.ID == productId {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	NotFound(w, r)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request, productId int) {
	var updateProduct Product
	err := json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		InvalidRequest(w, r)
		return
	}

	for i := range Products {
		if Products[i].ID == productId {
			updateProduct.ID = productId
			Products[i] = updateProduct
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduct)
			return
		}
	}

	NotFound(w, r)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request, productId int) {
	for i := range Products {
		if Products[i].ID == productId {
			Products = append(Products[:i], Products[i+1:]...)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
			return
		}
	}

	NotFound(w, r)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Categories)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		InvalidRequest(w, r)
		return
	}
	category.ID = len(Categories) + 1
	Categories = append(Categories, category)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func GetCategory(w http.ResponseWriter, r *http.Request, categoryId int) {
	for _, category := range Categories {
		if category.ID == categoryId {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	NotFound(w, r)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request, categoryId int) {
	var updateCategory Category
	err := json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		InvalidRequest(w, r)
		return
	}

	for i := range Categories {
		if Categories[i].ID == categoryId {
			updateCategory.ID = categoryId
			Categories[i] = updateCategory
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	NotFound(w, r)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request, categoryId int) {
	for i := range Categories {
		if Categories[i].ID == categoryId {
			Categories = append(Categories[:i], Categories[i+1:]...)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Category deleted"})
			return
		}
	}

	NotFound(w, r)
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.ReadInConfig()
	}

	config := Config{
		Port: viper.GetString("PORT"),
	}
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetProducts(w, r)
		case "POST":
			CreateProduct(w, r)
		default:
			UnknownMethod(w, r)
		}
	})

	http.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/products/")
		productId, err := strconv.Atoi(id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid product ID"})
			return
		}

		switch r.Method {
		case "GET":
			GetProduct(w, r, productId)
		case "PUT":
			UpdateProduct(w, r, productId)
		case "DELETE":
			DeleteProduct(w, r, productId)
		default:
			UnknownMethod(w, r)
		}
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetCategories(w, r)
		case "POST":
			CreateCategory(w, r)
		default:
			UnknownMethod(w, r)
		}
	})

	http.HandleFunc("/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/categories/")
		categoryId, err := strconv.Atoi(id)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid category ID"})
			return
		}

		switch r.Method {
		case "GET":
			GetCategory(w, r, categoryId)
		case "PUT":
			UpdateCategory(w, r, categoryId)
		case "DELETE":
			DeleteCategory(w, r, categoryId)
		default:
			UnknownMethod(w, r)
		}
	})

	fmt.Println("Server started on :" + config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
