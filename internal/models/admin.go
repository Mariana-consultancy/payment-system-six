package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	LockedAt     time.Time `json:"locked_at"`
	IsLocked     bool      `json:"is_locked"`
	FirstName    string    `json:"first_name" binding:"required"`
	Lastname     string    `json:"last_name" binding:"required"`
	Password     string    `json:"password" binding:"required"`
	DateOfBirth  string    `json:"date_of_birth" binding:"required"`
	Email        string    `json:"email" binding:"required,email"`
	Phone        string    `json:"phone" binding:"required"`
	Address      string    `json:"address" binding:"required"`
	LoginCounter int       `json:"login_counter"`
}

type LoginRequestAdmin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
