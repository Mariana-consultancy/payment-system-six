package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	LockedAt time.Time `json:"locked_at"`
	//IsLocked     bool      `json:"is_locked"`
	AccountNumber uint   `json:"account_number" gorm:"unique;not null"`
	FirstName     string `json:"first_name" binding:"required"`
	Lastname      string `json:"last_name" binding:"required"`
	Password      string `json:"password" binding:"required"`
	DateOfBirth   string `json:"date_of_birth" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone" binding:"required"`
	Address       string `json:"address" binding:"required"`
	//LoginCounter int       `json:"login_counter"`
	Balance float64 `json:"balance"`
}

type LoginRequestUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
