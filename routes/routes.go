package routes

import (
	"fmt"
	"net/http"

	"github.com/Pijuyy/testing_project4/controllers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "API Project 4 Kelompok 2")
	})
	// User routes
	router.HandleFunc("/users/register", controllers.RegisterUser(db)).Methods("POST")
	router.HandleFunc("/users/login", controllers.LoginUser(db)).Methods("POST")
	router.HandleFunc("/users/topup", controllers.TopUpUser(db)).Methods("PATCH")

	// Menggunakan instance CategoryController
	categoryController := controllers.NewCategoryController(db)

	// Category routes
	router.HandleFunc("/categories", categoryController.CreateCategory).Methods("POST")
	router.HandleFunc("/categories", categoryController.GetCategories).Methods("GET")
	router.HandleFunc("/categories/{categoryId}", categoryController.UpdateCategory).Methods("PATCH")
	router.HandleFunc("/categories/{categoryId}", categoryController.DeleteCategory).Methods("DELETE")

	// Product routes
	router.HandleFunc("/products", controllers.CreateProduct(db)).Methods("POST")
	router.HandleFunc("/products", controllers.GetProducts(db)).Methods("GET")
	router.HandleFunc("/products/{productId}", controllers.UpdateProduct(db)).Methods("PUT")
	router.HandleFunc("/products/{productId}", controllers.DeleteProduct(db)).Methods("DELETE")

	// TransactionHistory routes
	router.HandleFunc("/transactions", controllers.CreateTransaction(db)).Methods("POST")
	router.HandleFunc("/transactions/my-transactions", controllers.GetMyTransactions(db)).Methods("GET")
	router.HandleFunc("/transactions/user-transactions", controllers.GetUserTransactions(db)).Methods("GET")
}
