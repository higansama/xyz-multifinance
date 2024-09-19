package otp

// OTPServiceOpts contains options for OTPService.
type ServiceOpts struct {
}

// OTPService provides methods for generating and validating TOTP codes.
type Service struct {
}

// IOTPService is an interface for OTPService.
type IOTPService interface {
	// Generate generates a new TOTP code.
	//
	// Parameters:
	//    - opts: GenerateOTPOpts
	//       - secret: secret key
	//       - period: period in seconds
	//       - digits: number of digits
	//
	// Returns:
	//    - otpCode: otp code
	//    - countDown: countdown to the next period
	//    - error if any
	Generate(opts *GenerateOTPOpts) (string, uint, error)

	// Validate validates a TOTP code.
	//
	// Parameters:
	//    - otpCode: otp code
	//    - opts: GenerateOTPOpts (see Generate() for details)
	//
	// Returns:
	//    - valid: true if valid, false otherwise
	Validate(otpCode string, opts *GenerateOTPOpts) (bool, error)
}

// GenerateOTPOpts contains options for generating TOTP codes.
type GenerateOTPOpts struct {
	// Secret used to generate the TOTP hash.
	Secret string

	// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
	Period uint

	// Digits to request. Defaults to 6.
	Digits int
}
