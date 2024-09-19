package processor

import (
	"context"

	"github.com/go-playground/mold/v4"
	"github.com/higansama/xyz-multi-finance/internal/validator"
	"github.com/pkg/errors"
)

type StructProcessor struct {
	conforms  *mold.Transformer
	validator *validator.Validator
}

func NewStructProcessor(conforms *mold.Transformer, validator *validator.Validator) StructProcessor {
	return StructProcessor{
		conforms:  conforms,
		validator: validator,
	}
}

func (sp *StructProcessor) ProcessAndValidate(ctx context.Context, obj any) error {
	err := sp.conforms.Struct(ctx, obj)
	if err != nil {
		return errors.WithStack(err)
	}

	err = sp.validator.Validate(obj)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
