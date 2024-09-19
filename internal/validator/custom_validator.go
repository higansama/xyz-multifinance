package validator

import (
	"fmt"
	"reflect"

	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/higansama/xyz-multi-finance/internal/utils"
)

func RegisterCustomValidator(v *Validator) {
	if err := registerFilledValidation(v); err != nil {
		panic(err)
	}
	if err := registerUniqueStructFieldValidation(v); err != nil {
		panic(err)
	}
	if err := registerPhoneCountryCodeValidation(v); err != nil {
		panic(err)
	}

	reconstructParamInValidation(v, "gtfield")
	reconstructParamInValidation(v, "gtefield")
	reconstructParamInValidation(v, "eqfield")
	reconstructParamInValidation(v, "nefield")
	reconstructParamInValidation(v, "ltfield")
	reconstructParamInValidation(v, "ltefield")
}

func reconstructParamInValidation(v *Validator, tag string) {
	err := v.engine.RegisterTranslation(tag, v.translator, func(ut ut.Translator) error {
		return nil
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T(tag, fe.Field(), "::"+fe.Param()+"||::")
		return t
	})
	panicIf(err)
}

func registerFilledValidation(v *Validator) (err error) {
	// Field must be filled if it's present
	validationFunc := func(fl govalidator.FieldLevel) bool {
		field := fl.Field()

		switch field.Kind() {
		case reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan, reflect.Func:
			return !field.IsNil()
		case reflect.Ptr:
			return field.IsNil()
		default:
			if field.Interface() != nil {
				return field.Interface() != ""
			}
			return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
		}
	}

	err = v.engine.RegisterValidation("filled", validationFunc, true)
	if err != nil {
		return
	}

	err = v.engine.RegisterTranslation("filled", v.translator, func(ut ut.Translator) error {
		return ut.Add("filled", "{0} must be filled", false)
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T("filled", fe.Field())
		return t
	})
	if err != nil {
		return
	}

	return
}

func registerPhoneCountryCodeValidation(v *Validator) (err error) {
	fn := func(fl govalidator.FieldLevel) bool {
		field := fl.Field()

		switch field.Kind() {
		case reflect.String:
			val := field.String()
			if val == "" {
				return false
			}
			placeholder := "6981726"
			code := utils.GetCountryCodeFromPhoneNumber(val+placeholder, "")
			return code == field.String()
		default:
			panic(fmt.Sprintf("Bad field type %T", field.Interface()))
		}
	}

	err = makeValidValidation(v, fn, "phone_country_code")
	if err != nil {
		return
	}

	return
}

func registerUniqueStructFieldValidation(v *Validator) (err error) {
	validationFunc := func(fl govalidator.FieldLevel) bool {
		field := fl.Field()
		param := fl.Param()
		v := reflect.ValueOf(struct{}{})

		switch field.Kind() {
		case reflect.Slice, reflect.Array:
			elem := field.Type().Elem()
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}

			sf, ok := elem.FieldByName(param)
			if !ok {
				panic(fmt.Sprintf("Bad field name %s", param))
			}

			sfTyp := sf.Type
			if sfTyp.Kind() == reflect.Ptr {
				sfTyp = sfTyp.Elem()
			}

			m := reflect.MakeMap(reflect.MapOf(sfTyp, v.Type()))
			for i := 0; i < field.Len(); i++ {
				if field.Index(i).Kind() != reflect.Struct {
					panic("Not a struct type")
				}
				m.SetMapIndex(reflect.Indirect(reflect.Indirect(field.Index(i)).FieldByName(param)), v)
			}
			return field.Len() == m.Len()
		default:
			panic(fmt.Sprintf("Unsupported field type %T", field.Interface()))
		}
	}

	err = v.engine.RegisterValidation("unique_sf", validationFunc, true)
	if err != nil {
		return
	}

	err = v.engine.RegisterTranslation("unique_sf", v.translator, func(ut ut.Translator) error {
		return ut.Add("unique_sf", "{0} ({1}) should be unique", true)
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T("unique_sf", fe.Field(), "::"+fe.Param()+"||::")
		return t
	})
	if err != nil {
		return
	}

	return
}

// The field under validation must be present and not empty only if any the other specified fields
// are equal to the value following with the specified field.
func registerRequiredIfAnyValidation(v *Validator) (err error) {
	fn := func(fl govalidator.FieldLevel) bool {
		//field := fl.Field()
		//param := fl.Param()

		otherField, otherKind, _, ok := fl.GetStructFieldOK2()
		if !ok {
			panic("Template Kind param is required for Template Name")
		}

		if otherKind != reflect.String {
			panic(fmt.Sprintf("Bad field type %T", otherField.Interface()))
		}

		//

		return true
	}

	err = v.engine.RegisterValidation("required_if_any", fn, true)
	if err != nil {
		return
	}

	err = v.engine.RegisterTranslation("required_if_any", v.translator, func(ut ut.Translator) error {
		return ut.Add("required_if_any", "{0} ({1}) is a required field", true)
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T("unique_sf", fe.Field(), "::"+fe.Param()+"||::")
		return t
	})
	if err != nil {
		return
	}

	return
}
