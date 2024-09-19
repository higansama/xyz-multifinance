package auth

import (
	// "config"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"github.com/higansama/xyz-multi-finance/config"
	"github.com/pkg/errors"
)

type JwtClaims struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
	// UserIsAMaster    bool   `json:"user_master"`
	// UserIsADoa       bool   `json:"user_doa"`
	// OldUserId        string `json:"_id"`
	// OldUserIsAMaster bool `json:"master"`
	jwt.RegisteredClaims
}

type GenerateAuthTokenOptions struct {
	UserId    string
	UserName  string
	UserEmail string
	UserRole  string
}

func GenerateAuthToken(config config.Config, opts GenerateAuthTokenOptions) (string, error) {
	claims := JwtClaims{
		UserId:           opts.UserId,
		UserName:         opts.UserName,
		UserEmail:        opts.UserEmail,
		UserRole:         opts.UserRole,
		RegisteredClaims: jwt.RegisteredClaims{IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(carbon.Now().AddDay().ToStdTime())},
	}

	jwtSigner := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtSigner.SignedString([]byte(config.Auth.JwtSecret))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}
