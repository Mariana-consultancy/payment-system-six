package repository

import (
	"fmt"
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"strconv"

	"gorm.io/gorm"
)

func (p *Postgres) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (p *Postgres) FindAllUsers() ([]models.User, error) {
	var users []models.User

	if err := p.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (p *Postgres) GetUserByID(userID uint) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("ID = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (p *Postgres) GetUserByAccountNumber(accountNumber uint) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("account_number = ?", accountNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Create a user in the database
func (p *Postgres) CreateUser(user *models.User) error {
	if err := p.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// Update a user in the database
func (p *Postgres) UpdateUser(user *models.User) error {
	if err := p.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// Generate Account Number For User
func (p *Postgres) GenerateUserAccountNumber() (uint, error) {
	for {
		number, err := util.Generate8DigitNumber()
		if err != nil {
			return 0, err
		}
		user := &models.User{}
		if err := p.DB.Where("account_number = ?", number).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Number is unique
				return number, nil
			}
			// Some other error occurred
			return 0, err
		}
		// Number already exists, generate a new one
	}
}

func (p *Postgres) GenerateStatement(accountNumber uint) (*models.StatementDetails, error) {
	statementDetails := &models.StatementDetails{}
	transactions := &[]models.Transaction{}
	if err := p.DB.Where("account_number = ?", accountNumber).Order("created_at").Find(&transactions).Error; err != nil {
		return nil, err
	}

	if len(*transactions) == 0 {
		return nil, fmt.Errorf("no transactions found for account number %d", accountNumber)
	}

	// Calculate the opening balance from the first transaction
	statementDetails.OpeningBalance = (*transactions)[0].BalanceBefore

	var totalIn float64
	var totalOut float64

	for _, transaction := range *transactions {
		statementEntry := models.Statement{
			Date:           transaction.TransactionDate,
			Type:           transaction.TransactionType,
			AccountBalance: transaction.BalanceAfter,
		}

		if transaction.PayerAccountNumber == accountNumber && transaction.PayeeAccountNumber == accountNumber {
			statementEntry.Transaction = "Deposit"
			statementEntry.In = transaction.TransactionAmount
			totalIn += transaction.TransactionAmount
		} else if transaction.PayerAccountNumber == accountNumber {
			// This is an outgoing transaction
			statementEntry.Transaction = "Transfer to Acc#  " + strconv.FormatUint(uint64(transaction.PayeeAccountNumber), 10)
			statementEntry.Out = transaction.TransactionAmount
			totalOut += transaction.TransactionAmount

		} else if transaction.PayeeAccountNumber == accountNumber {
			// This is an incoming transaction
			statementEntry.Transaction = "Received from Acc#  " + strconv.FormatUint(uint64(transaction.PayerAccountNumber), 10)
			statementEntry.In = transaction.TransactionAmount
			totalIn += transaction.TransactionAmount
		}

		statementDetails.Statement = append(statementDetails.Statement, statementEntry)
	}

	statementDetails.TotalIn = totalIn
	statementDetails.TotalOut = totalOut
	statementDetails.ClosingBlance = (*transactions)[len(*transactions)-1].BalanceAfter

	return statementDetails, nil
}
