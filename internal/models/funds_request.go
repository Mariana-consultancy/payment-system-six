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
	CardDetails
}

type MakePayment struct {
	PayeeEmail string  `json:"payee_email" binding:"required,email"`
	Amount     float64 `json:"amount" binding:"required"`
}

type PaymentRequests struct {
	gorm.Model
	RequesterID  uint    `json:"requester_id"`
	RecipientID  uint    `json:"recipient_id"`
	Amount       float64 `json:"transaction_amount"`
	Status       string  `json:"balance_before"`
	Desscription string  `json:"description"`
}

type RequestPayment struct {
	RecipientEmail string  `json:"recipient_email" binding:"required,email"`
	Amount         float64 `json:"amount" binding:"required"`
	Desscription   string  `json:"description"`
}
