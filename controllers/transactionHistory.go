package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Pijuyy/testing_project4/config"
	"github.com/Pijuyy/testing_project4/models"
	"gorm.io/gorm"
)

// CreateTransaction - Create a new transaction
func CreateTransaction(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate user
		user, err := config.ExtractUserFromToken(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Parse request body
		var requestBody struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the specified product exists
		var product models.Product
		result := db.First(&product, requestBody.ProductID)
		if result.Error != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Check if the quantity is available in stock
		if requestBody.Quantity > product.Stock {
			http.Error(w, "Not enough stock available", http.StatusBadRequest)
			return
		}

		// Check if user has enough balance
		totalPrice := product.Price * requestBody.Quantity
		if user.Balance < int64(totalPrice) {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		// Deduct stock from the product
		product.Stock -= requestBody.Quantity
		db.Save(&product)

		// Deduct balance from the user
		user.Balance -= int64(totalPrice)
		db.Save(&user)

		// Update sold_product_amount in category
		var category models.Category
		db.First(&category, product.CategoryID)
		category.SoldProductAmount += requestBody.Quantity
		db.Save(&category)

		// Create a new transaction history record
		transactionHistory := models.TransactionHistory{
			ProductID:  requestBody.ProductID,
			UserID:     user.ID,
			Quantity:   requestBody.Quantity,
			TotalPrice: totalPrice,
			CreatedAt:  time.Now(),
		}
		db.Create(&transactionHistory)

		// Prepare the response
		response := map[string]interface{}{
			"message": "You have successfully purchased the product",
			"transaction_bill": map[string]interface{}{
				"total_price":   totalPrice,
				"quantity":      requestBody.Quantity,
				"product_title": product.Title,
			},
		}

		// Send the JSON response
		w.WriteHeader(http.StatusCreated)
		config.SendJSONResponse(w, response)
	}
}

// GetMyTransactions - Get transactions of the authenticated user
func GetMyTransactions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate user
		user, err := config.ExtractUserFromToken(r, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Retrieve user's transactions
		var transactions []models.TransactionHistory
		db.Preload("Product").Where("user_id = ?", user.ID).Find(&transactions)

		// Prepare the response
		response := make([]map[string]interface{}, 0)
		for _, transaction := range transactions {
			transactionData := map[string]interface{}{
				"id":          transaction.ID,
				"product_id":  transaction.ProductID,
				"user_id":     transaction.UserID,
				"quantity":    transaction.Quantity,
				"total_price": transaction.TotalPrice,
				"Product": map[string]interface{}{
					"id":          transaction.Product.ID,
					"title":       transaction.Product.Title,
					"price":       transaction.Product.Price,
					"stock":       transaction.Product.Stock,
					"category_id": transaction.Product.CategoryID,
					"created_at":  transaction.Product.CreatedAt,
					"updated_at":  transaction.Product.UpdatedAt,
				},
			}
			response = append(response, transactionData)
		}

		// Send the JSON response
		config.SendJSONResponse(w, response)
	}
}

// GetUserTransactions - Get all transactions for admin
func GetUserTransactions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and authorize user
		authorized, err := config.AuthenticateAndAuthorize(r, db)
		if err != nil || !authorized {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		// Retrieve all transactions
		var transactions []models.TransactionHistory
		db.Preload("Product").Preload("User").Find(&transactions)

		// Prepare the response
		response := make([]map[string]interface{}, 0)
		for _, transaction := range transactions {
			transactionData := map[string]interface{}{
				"id":          transaction.ID,
				"product_id":  transaction.ProductID,
				"user_id":     transaction.UserID,
				"quantity":    transaction.Quantity,
				"total_price": transaction.TotalPrice,
				"Product": map[string]interface{}{
					"id":          transaction.Product.ID,
					"title":       transaction.Product.Title,
					"price":       transaction.Product.Price,
					"stock":       transaction.Product.Stock,
					"category_id": transaction.Product.CategoryID,
					"created_at":  transaction.Product.CreatedAt,
					"updated_at":  transaction.Product.UpdatedAt,
				},
				"User": map[string]interface{}{
					"id":         transaction.User.ID,
					"email":      transaction.User.Email,
					"full_name":  transaction.User.FullName,
					"balance":    transaction.User.Balance,
					"created_at": transaction.User.CreatedAt,
					"updated_at": transaction.User.UpdatedAt,
				},
			}
			response = append(response, transactionData)
		}

		// Send the JSON response
		config.SendJSONResponse(w, response)
	}
}
