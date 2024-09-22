package usecasev1

import "github.com/higansama/xyz-multi-finance/customer/dto"

// GetMyLimit implements customer.Usecase.
func (c *CustomerUsecase) GetMyLimit(userid string) (response dto.MyLimitResponse, err error) {
	r, err := c.repo.GetUser(userid)
	if err != nil {
		return response, err
	}

	response = dto.MyLimitResponse{
		Userid:     userid,
		Email:      r.Email,
		OneMonth:   int(r.MyLoanLimit().OneMonth),
		TwoMonth:   int(r.MyLoanLimit().TwoMonth),
		ThreeMonth: int(r.MyLoanLimit().ThreeMonth),
		FourMonth:  int(r.MyLoanLimit().FourMonth),
	}
	return response, nil
}
