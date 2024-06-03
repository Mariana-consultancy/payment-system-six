package repository

import "payment-system-six/internal/models"

func (p *Postgres) FindAdminByEmail(email string) (*models.Admin, error) {
	admin := &models.Admin{}

	if err := p.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

// Create a user in the database
func (p *Postgres) CreateAdmin(admin *models.Admin) error {
	if err := p.DB.Create(admin).Error; err != nil {
		return err
	}

	return nil
}

// Update a user in the database
func (p *Postgres) UpdateAdmin(admin *models.Admin) error {
	if err := p.DB.Save(admin).Error; err != nil {
		return err
	}
	return nil
}
