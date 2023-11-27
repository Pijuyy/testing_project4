package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID         uint   `gorm:"primary_key"`
	Title      string `gorm:"not null"`
	Price      int    `gorm:"not null"`
	Stock      int    `gorm:"not null"`
	CategoryID uint
	Category   Category `gorm:"foreignKey:CategoryID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	// Validate title
	if p.Title == "" {
		return errors.New("title is required")
	}

	// Validate stock
	if p.Stock < 5 {
		return errors.New("stock must be at least 5")
	}

	// Validate price
	if p.Price < 0 || p.Price > 50000000 {
		return errors.New("price must be between 0 and 50,000,000")
	}

	return
}
