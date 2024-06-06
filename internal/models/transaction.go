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

type StatementDetails struct {
	OpeningBalance float64     `json:"opening_balance"`
	TotalIn        float64     `json:"total_in"`
	TotalOut       float64     `json:"total_out"`
	ClosingBlance  float64     `json:"closing_balance"`
	Statement      []Statement `json:"statement"`
}

type Statement struct {
	Date           time.Time `json:"date"`
	Type           string    `json:"type"`
	Transaction    string    `json:"transaction"`
	In             float64   `json:"in"`
	Out            float64   `json:"out"`
	AccountBalance float64   `json:"account_balance"`
}
