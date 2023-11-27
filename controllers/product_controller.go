package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Pijuyy/testing_project4/config"
	"github.com/Pijuyy/testing_project4/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateProduct - Create a new product
func CreateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		_, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var requestBody struct {
			Title      string `json:"title"`
			Price      int    `json:"price"`
			Stock      int    `json:"stock"`
			CategoryID int    `json:"category_id"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the specified category exists
		var category models.Category
		result := db.First(&category, requestBody.CategoryID)
		if result.Error != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		product := models.Product{
			Title:      requestBody.Title,
			Price:      requestBody.Price,
			Stock:      requestBody.Stock,
			CategoryID: uint(requestBody.CategoryID),
			CreatedAt:  time.Now(),
		}

		result = db.Create(&product)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		config.SendJSONResponse(w, map[string]interface{}{
			"id":          product.ID,
			"title":       product.Title,
			"price":       product.Price,
			"stock":       product.Stock,
			"category_id": product.CategoryID,
			"created_at":  product.CreatedAt.Format(time.RFC3339),
		})
	}
}

// GetProducts - Get all products
func GetProducts(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		_, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var products []models.Product
		result := db.Find(&products)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Create a structured response with ordered fields
		var responseProducts []map[string]interface{}
		for _, product := range products {
			productData := map[string]interface{}{
				"id":          product.ID,
				"title":       product.Title,
				"price":       product.Price,
				"stock":       product.Stock,
				"category_id": product.CategoryID,
				"created_at":  product.CreatedAt.Format(time.RFC3339),
			}

			responseProducts = append(responseProducts, productData)
		}

		// Send the JSON response with ordered fields
		config.SendJSONResponse(w, responseProducts)
	}
}

// UpdateProduct - Update product by ID
func UpdateProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		_, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		productID, err := strconv.Atoi(mux.Vars(r)["productId"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var requestBody struct {
			Title      string `json:"title"`
			Price      int    `json:"price"`
			Stock      int    `json:"stock"`
			CategoryID int    `json:"category_id"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the specified category exists
		var category models.Category
		result := db.First(&category, requestBody.CategoryID)
		if result.Error != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		var product models.Product
		result = db.First(&product, productID)
		if result.Error != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		product.Title = requestBody.Title
		product.Price = requestBody.Price
		product.Stock = requestBody.Stock
		product.CategoryID = uint(requestBody.CategoryID)
		product.UpdatedAt = time.Now()

		db.Save(&product)

		config.SendJSONResponse(w, map[string]interface{}{
			"Products": map[string]interface{}{
				"id":          product.ID,
				"title":       product.Title,
				"price":       product.Price,
				"stock":       product.Stock,
				"category_id": product.CategoryID,
				"created_at":  product.CreatedAt.Format(time.RFC3339),
				"updated_at":  product.UpdatedAt.Format(time.RFC3339),
			},
		})
	}
}

// DeleteProduct - Delete product by ID
func DeleteProduct(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize the user
		_, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		productID, err := strconv.Atoi(mux.Vars(r)["productId"])
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var product models.Product
		result := db.First(&product, productID)
		if result.Error != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		db.Delete(&product)

		config.SendJSONResponse(w, map[string]interface{}{
			"message": "Product has been successfully deleted",
		})
	}
}
