package config

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Pijuyy/testing_project4/models"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func init() {
	if len(JwtKey) == 0 {
		JwtKey = []byte("bMLRrApp4zf6qzWoMa-brT6HMwG5Lp5VY8l1Y-K34Xwsm8B3-kB9p7pcRoWKP8jafaTuCylxPMllgz6uFT6zfQ") // Fallback key
	}
}

func SendJSONResponse(w http.ResponseWriter, v interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(v)
}

func ExtractUserFromToken(r *http.Request, db *gorm.DB) (*models.User, error) {
	claims, err := Authenticate(r, db)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func Authenticate(r *http.Request, db *gorm.DB) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return nil, errors.New("Authorization header must be in the format 'Bearer {token}'")
	}
	tokenString := authHeaderParts[1]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		return nil, errors.New("Invalid token")
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	// Fetch user based on the email from claims
	var user models.User
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		return nil, errors.New("User not found")
	}

	// Set the correct UserID in the claims
	claims.UserID = user.ID

	return claims, nil
}

func AuthenticateAndAuthorize(r *http.Request, db *gorm.DB) (bool, error) {
	claims, err := Authenticate(r, db)
	if err != nil {
		return false, err
	}

	var user models.User
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		return false, errors.New("User not found")
	}

	if user.Role != "admin" {
		return false, errors.New("Unauthorized access")
	}

	return true, nil
}
