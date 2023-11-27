package models

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	FullName  string `gorm:"not null" valid:"required"`
	Email     string `gorm:"not null;unique" valid:"email,required"`
	Password  string `gorm:"not null" valid:"required,length(6|255)"`
	Role      string `gorm:"not null" valid:"required,oneof=admin customer" json:"Role"`
	Balance   int64  `gorm:"not null" valid:"range(0|100000000)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Validate email
	if !govalidator.IsEmail(u.Email) {
		return errors.New("invalid email")
	}

	// Validate full name
	if u.FullName == "" {
		return errors.New("full name is required")
	}

	// Validate balance
	if u.Balance < 0 || u.Balance > 100000000 {
		return errors.New("balance must be between 0 and 100,000,000")
	}

	// Validate password
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	// Validate role
	if u.Role != "admin" && u.Role != "customer" {
		return errors.New("role must be either 'admin' or 'customer'")
	}

	return
}
