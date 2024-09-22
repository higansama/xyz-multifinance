package usecasev1

import (
	"time"

	domain "github.com/higansama/xyz-multi-finance/customer"
	"github.com/higansama/xyz-multi-finance/customer/dto"
	"github.com/higansama/xyz-multi-finance/customer/errors"
)

// Create implements customer.Usecase.
func (c *CustomerUsecase) Create(data dto.ReqCreateCustomer) error {
	// find user

	user, err := c.repo.GetUser(data.Email)
	if err != nil {
		return err
	}

	if string(user.ID) != "" {
		return errors.EmailAtauNIKTerdaftar
	}

	// email must be unique
	customer := domain.CustomerEntity{
		NIK:        data.NIK,
		Email:      data.Email,
		FullName:   data.FullName,
		LegalName:  data.LegalName,
		BirthPlace: data.BirthPlace,
		BirthDate:  data.BirthDate,
		Salary:     data.Salary,
		PhotoID:    data.PhotoID,
		Selfie:     data.Selfie,
		CreatedAt:  time.Now(),
	}
	// hash the password
	customer.HashPassword(data.Password)
	err = c.repo.Create(customer)
	if err != nil {
		return err
	}
	return nil
}
