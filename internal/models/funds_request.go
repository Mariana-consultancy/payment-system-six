package models

import (
	"gorm.io/gorm"
)

type CardDetails struct {
	CardNumber string `json:"card_number" binding:"required"`
	CardExpiry string `json:"card_expiry" binding:"required"`
	CardCVV    string `json:"card_cvv" binding:"required"`
}

type FundRequest struct {
	Amount float64 `json:"amount" binding:"required"`
	//CardDetails
}

type TransferPayment struct {
	AccountNumber uint    `json:"account_number" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
}

type PaymentRequests struct {
	gorm.Model
	RequesterAccountNumber uint    `json:"requester_account_number"`
	RecipientAccountNumber uint    `json:"recipient_account_number"`
	Amount                 float64 `json:"transaction_amount"`
	Status                 string  `json:"balance_before"`
	Desscription           string  `json:"description"`
}

type RequestPayment struct {
	AccountNumber uint    `json:"account_number" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Desscription  string  `json:"description"`
}
