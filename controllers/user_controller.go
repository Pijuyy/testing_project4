package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Pijuyy/testing_project4/config"
	"github.com/Pijuyy/testing_project4/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser - Register User as a Customer
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			FullName string `json:"full_name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error while hashing password", http.StatusInternalServerError)
			return
		}

		// User role is hardcoded as 'customer' and balance starts from 0
		user := models.User{
			FullName: requestBody.FullName,
			Email:    requestBody.Email,
			Password: string(hashedPassword),
			Role:     "customer",
			Balance:  0,
		}

		result := db.Create(&user)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		config.SendJSONResponse(w, map[string]interface{}{
			"id":         user.ID,
			"full_name":  user.FullName,
			"email":      user.Email,
			"password":   user.Password, // Note: Avoid sending sensitive information in response
			"balance":    user.Balance,
			"created_at": user.CreatedAt.Format(time.RFC3339),
		})
	}
}

func InitializeAdminUser(db *gorm.DB) {
	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	adminUser := models.User{
		FullName: "Admin User",
		Email:    "admin@gmail.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	// Check if admin exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", adminUser.Email).Count(&count)
	if count == 0 {
		// Create admin if not exists
		db.Create(&adminUser)
	}
}

// LoginUser - Login User
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user models.User
		result := db.Where("email = ?", requestBody.Email).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(1 * time.Hour)
		claims := &config.Claims{
			Email: user.Email,
			Role:  user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(config.JwtKey)
		if err != nil {
			http.Error(w, "Error while signing the token", http.StatusInternalServerError)
			return
		}

		config.SendJSONResponse(w, map[string]string{
			"token": tokenString,
		})
	}
}

// TopupUserBalance - Top-up user balance
func TopUpUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate the user
		user, err := config.ExtractUserFromToken(r, db)
		if err != nil {
			http.Error(w, "Failed to extract user from token", http.StatusUnauthorized)
			return
		}

		var requestBody struct {
			Balance int64 `json:"balance"`
		}
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update user balance
		user.Balance += requestBody.Balance
		db.Save(&user)

		// Send a JSON response indicating successful top-up
		config.SendJSONResponse(w, map[string]string{
			"message": "Your balance has been successfully updated to Rp " + fmt.Sprintf("%d", user.Balance),
		})
	}
}
