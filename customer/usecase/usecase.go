package usecase

import (
	domain "github.com/higansama/xyz-multi-finance/customer"
	"github.com/higansama/xyz-multi-finance/internal/infrastructure"
)

type CustomerUsecase struct {
	infra *infrastructure.Infrastructure
	repo  domain.Repository
}

func NewUserUsecase(
	infra *infrastructure.Infrastructure,
	repo domain.Repository,
) (domain.Usecase, error) {
	u := &CustomerUsecase{
		infra: infra,
		repo:  repo,
	}

	return u, nil
}
