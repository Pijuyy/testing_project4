// service.go
package service

import (
	"github.com/Pijuyy/testing_project4/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		DB: db,
	}
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.DB.Create(category).Error
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	result := s.DB.Preload("Products").Find(&categories)
	return categories, result.Error
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	return s.DB.Save(category).Error
}

func (s *CategoryService) DeleteCategory(category *models.Category) error {
	return s.DB.Delete(category).Error
}
