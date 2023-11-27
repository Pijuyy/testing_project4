// repository.go
package repo

import (
	"github.com/Pijuyy/testing_project4/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		DB: db,
	}
}

func (r *CategoryRepository) FindCategoryByID(id int) (*models.Category, error) {
	var category models.Category
	result := r.DB.First(&category, id)
	return &category, result.Error
}
