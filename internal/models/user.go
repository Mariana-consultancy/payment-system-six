package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string    `json:"first_name"`
	Lastname     string    `json:"last_name"`
	Password     string    `json:"password"`
	DateOfBirth  string    `json:"date_of_birth"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	LoginCounter int       `json:"login_counter"`
	IsLocked     bool      `json:"is_locked"`
	LockedAt     time.Time `json:"locked_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `son:"password"`
}
