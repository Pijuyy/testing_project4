package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Category struct represents a category in the system.
type Category struct {
	ID                uint   `gorm:"primary_key"`
	Type              string `gorm:"not null"`
	SoldProductAmount int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Products          []Product `gorm:"foreignKey:CategoryID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {

	if c.Type == "" {
		return errors.New("type is required")
	}

	return
}
