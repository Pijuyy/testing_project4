// controller.go
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Pijuyy/testing_project4/config"
	"github.com/Pijuyy/testing_project4/models"
	"github.com/Pijuyy/testing_project4/repo"
	"github.com/Pijuyy/testing_project4/service"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CategoryController struct {
	Service    *service.CategoryService
	Repository *repo.CategoryRepository
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{
		Service:    service.NewCategoryService(db),
		Repository: repo.NewCategoryRepository(db),
	}
}

// CreateCategory - Create a new category
func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	user, err := config.ExtractUserFromToken(r, c.Service.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	var requestBody struct {
		Type string `json:"type"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category := models.Category{
		Type:              requestBody.Type,
		SoldProductAmount: 0,
		CreatedAt:         time.Now(),
	}

	err = c.Service.CreateCategory(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	config.SendJSONResponse(w, map[string]interface{}{
		"id":                  category.ID,
		"type":                category.Type,
		"sold_product_amount": category.SoldProductAmount,
		"created_at":          category.CreatedAt.Format(time.RFC3339),
	})
}

// GetCategories - Get all categories
func (c *CategoryController) GetCategories(w http.ResponseWriter, r *http.Request) {
	user, err := config.ExtractUserFromToken(r, c.Service.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	categories, err := c.Service.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a structured response with ordered fields
	var responseCategories []map[string]interface{}
	for _, category := range categories {
		categoryData := map[string]interface{}{
			"id":                  category.ID,
			"type":                category.Type,
			"sold_product_amount": category.SoldProductAmount,
			"created_at":          category.CreatedAt.Format(time.RFC3339),
			"updated_at":          category.UpdatedAt.Format(time.RFC3339),
			"products":            []map[string]interface{}{},
		}

		// Add product details to the response with ordered fields
		for _, product := range category.Products {
			productData := map[string]interface{}{
				"id":         product.ID,
				"title":      product.Title,
				"price":      product.Price,
				"stock":      product.Stock,
				"created_at": product.CreatedAt.Format(time.RFC3339),
				"updated_at": product.UpdatedAt.Format(time.RFC3339),
			}
			categoryData["products"] = append(categoryData["products"].([]map[string]interface{}), productData)
		}

		responseCategories = append(responseCategories, categoryData)
	}

	// Send the JSON response with ordered fields
	config.SendJSONResponse(w, responseCategories)
}

// UpdateCategory - Update category by ID
func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	user, err := config.ExtractUserFromToken(r, c.Service.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	categoryID, err := strconv.Atoi(mux.Vars(r)["categoryId"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Type string `json:"type"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := c.Repository.FindCategoryByID(categoryID)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	category.Type = requestBody.Type
	category.UpdatedAt = time.Now()

	err = c.Service.UpdateCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config.SendJSONResponse(w, map[string]interface{}{
		"id":                  category.ID,
		"type":                category.Type,
		"sold_product_amount": category.SoldProductAmount,
		"updated_at":          category.UpdatedAt.Format(time.RFC3339),
	})
}

// DeleteCategory - Delete category by ID
func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	user, err := config.ExtractUserFromToken(r, c.Service.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	categoryID, err := strconv.Atoi(mux.Vars(r)["categoryId"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := c.Repository.FindCategoryByID(categoryID)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	err = c.Service.DeleteCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config.SendJSONResponse(w, map[string]interface{}{
		"message": "Category has been successfully deleted",
	})
}
