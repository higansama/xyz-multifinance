package app

import (
	"github.com/pkg/errors"
)

type Environment struct {
	v string
	s string
}

var (
	EnvironmentLocal       = Environment{"local", "local"}
	EnvironmentDevelopment = Environment{"development", "dev"}
	EnvironmentTesting     = Environment{"testing", "test"}
	EnvironmentStaging     = Environment{"staging", "stag"}
	EnvironmentProduction  = Environment{"production", "prod"}
)

var EnvironmentList = []Environment{
	EnvironmentLocal,
	EnvironmentDevelopment,
	EnvironmentTesting,
	EnvironmentStaging,
	EnvironmentProduction,
}

func NewEnvironmentFromString(str string) (Environment, error) {
	for _, v := range EnvironmentList {
		if v.String() == str {
			return v, nil
		}
	}
	return Environment{}, errors.Errorf("unknown '%s' environment", str)
}

func (u Environment) IsZero() bool {
	return u == Environment{}
}

func (u Environment) Short() string {
	return u.s
}

func (u Environment) String() string {
	return u.v
}
