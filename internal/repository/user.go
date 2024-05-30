package repository

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"

	"gorm.io/gorm"
)

func (p *Postgres) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
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
