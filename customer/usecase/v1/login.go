package usecasev1

import (
	"github.com/higansama/xyz-multi-finance/internal/auth"

	"github.com/higansama/xyz-multi-finance/customer/dto"
	"github.com/higansama/xyz-multi-finance/customer/errors"
)

// Login implements customer.Usecase.
func (c *CustomerUsecase) Login(payload dto.LoginRequest) (response dto.LoginResponse, err error) {
	user, err := c.repo.GetUserByEmail(payload.Email)
	if err != nil {
		return response, err
	}

	if user.ID == "" {
		return response, errors.UserNotFound
	}

	err = user.ValidatePassword(payload.Password)
	if err != nil {
		return response, err
	}

	// create jwt
	jwtToken, err := auth.GenerateAuthToken(c.infra.Config, auth.GenerateAuthTokenOptions{
		UserId:    user.ID,
		UserName:  user.FullName,
		UserEmail: user.Email,
		Salary:    user.Salary,
	})

	if err != nil {
		return response, err
	}

	response.Token = jwtToken
	return response, nil
}
