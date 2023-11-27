package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type TransactionHistory struct {
	ID         uint `gorm:"primary_key"`
	ProductID  uint
	Product    Product `gorm:"foreignKey:ProductID"`
	UserID     uint
	User       User `gorm:"foreignKey:UserID"`
	Quantity   int  `gorm:"not null"`
	TotalPrice int  `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (th *TransactionHistory) BeforeCreate(tx *gorm.DB) (err error) {
	// Validate quantity
	if th.Quantity <= 0 {
		return errors.New("quantity is required and must be greater than 0")
	}

	// Validate total price
	if th.TotalPrice <= 0 {
		return errors.New("total price is required and must be greater than 0")
	}

	return
}
