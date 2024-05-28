package repository

import (
	"errors"
	"payment-system-six/internal/models"
)

func (p *Postgres) ValidateCard(cardNumber, cardExpiry, cardCVV string) error {

	if cardNumber != "1234 5678 9101 1121" || cardExpiry != "06/25" || cardCVV != "123" {
		return errors.New("invalid card details")
	}
	return nil
}

func (p *Postgres) RecordTransaction(transaction *models.Transaction) error {
	if err := p.DB.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) RequestFunds(paymentRequest *models.PaymentRequests) error {
	if err := p.DB.Create(paymentRequest).Error; err != nil {
		return err
	}
	return nil
}
