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

func (p *Postgres) GetAllPaymentRequestsByAccountNumber(accountNumber uint) (*[]models.PaymentRequests, error) {
	paymentRequests := &[]models.PaymentRequests{}
	if err := p.DB.Where("recipient_account_number = ? OR requester_account_number = ?", accountNumber, accountNumber).
		Order("CASE WHEN status = 'pending' THEN 0 ELSE 1 END").
		Order("updated_at DESC").
		Find(&paymentRequests).Error; err != nil {
		return nil, err
	}

	return paymentRequests, nil
}

func (p *Postgres) GetPaymentRequestByRequestID(requestID uint) (*models.PaymentRequests, error) {
	paymentRequest := &models.PaymentRequests{}
	if err := p.DB.Where("ID = ?", requestID).First(&paymentRequest).Error; err != nil {
		return nil, err
	}
	return paymentRequest, nil
}

func (p *Postgres) UpdatePaymentRequest(paymentRequest *models.PaymentRequests) error {
	if err := p.DB.Save(paymentRequest).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DeletePaymentRequest(paymentRequest *models.PaymentRequests) error {
	if err := p.DB.Delete(paymentRequest).Error; err != nil {
		return err
	}
	return nil
}
