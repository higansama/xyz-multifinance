package validator

import (
	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
)

func RegisterValidatorTranslations(v *Validator) {
	err := v.engine.RegisterTranslation("http_url", v.translator, func(ut ut.Translator) error {
		return ut.Add("http_url", "{0} is not a valid URL", false)
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T("http_url", fe.Field())
		return t
	})
	panicIf(err)
}
