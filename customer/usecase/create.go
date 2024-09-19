package usecase

import domain "github.com/higansama/xyz-multi-finance/customer"

// Create implements customer.Usecase.
func (c *CustomerUsecase) Create(data domain.CustomerEntity) error {
	c.repo.GetUser("nashrulloh1811@gmail.com")
	return nil
}
