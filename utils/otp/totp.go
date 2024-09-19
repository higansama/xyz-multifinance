package otp

import (
	"encoding/base32"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func NewOTPService(opts *ServiceOpts) IOTPService {
	return &Service{}
}

func (o *Service) Generate(opts *GenerateOTPOpts) (string, uint, error) {

	// create a new secret each time
	secret := base32.StdEncoding.EncodeToString([]byte(opts.Secret))

	// count down to the next period
	currentTime := time.Now()
	countDown := opts.Period - uint(currentTime.Unix()%int64(opts.Period))

	// generate otp code
	otpCode, err := totp.GenerateCodeCustom(secret, currentTime, totp.ValidateOpts{
		Period: opts.Period,
		Digits: otp.Digits(opts.Digits),
	})

	if err != nil {
		return otpCode, 0, err
	}

	return otpCode, countDown, nil
}

func (o *Service) Validate(otpCode string, opts *GenerateOTPOpts) (bool, error) {
	// create a new secret each time
	secret := base32.StdEncoding.EncodeToString([]byte(opts.Secret))

	// validate otp code
	valid, err := totp.ValidateCustom(otpCode, secret, time.Now(), totp.ValidateOpts{
		Period: opts.Period,
		Digits: otp.Digits(opts.Digits),
	})

	if err != nil {
		return false, err
	}

	return valid, nil
}
