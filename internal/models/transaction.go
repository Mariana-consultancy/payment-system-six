package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	PayerAccountNumber uint      `json:"payer_account_number"`
	PayeeAccountNumber uint      `json:"payee_account_number"`
	AccountNumber      uint      `json:"account_number"`
	TransactionType    string    `json:"transaction_type"`
	TransactionAmount  float64   `json:"transaction_amount"`
	BalanceBefore      float64   `json:"balance_before"`
	BalanceAfter       float64   `json:"balance_after"`
	TransactionDate    time.Time `json:"transaction_date"`
}
