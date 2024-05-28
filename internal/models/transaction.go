package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	PayerID           uint      `json:"payer_id"`
	PayeeID           uint      `json:"payee_id"`
	AccountID         uint      `json:"account_id"`
	TransactionType   string    `json:"transaction_type"`
	TransactionAmount float64   `json:"transaction_amount"`
	BalanceBefore     float64   `json:"balance_before"`
	BalanceAfter      float64   `json:"balance_after"`
	TransactionDate   time.Time `json:"transaction_date"`
}
