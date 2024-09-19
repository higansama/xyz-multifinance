package validator

// func RegisterSpecificValidator(v *Validator) {
// 	if err := registerTemplateNameValidation(v); err != nil {
// 		panic(err)
// 	}
// 	if err := registerTemplateCodeIdValidation(v); err != nil {
// 		panic(err)
// 	}
// }

// func registerTemplateNameValidation(v *Validator) (err error) {
// 	fn := func(fl govalidator.FieldLevel) bool {
// 		field := fl.Field()

// 		otherField, otherKind, _, ok := fl.GetStructFieldOK2()
// 		if !ok {
// 			panic("Template Kind param is required for Template Name")
// 		}

// 		if otherKind != reflect.String {
// 			panic(fmt.Sprintf("Bad field type %T", otherField.Interface()))
// 		}

// 		otherValue := otherField.String()
// 		if templateVo.TemplateKindWhatsapp.Equals(otherValue) ||
// 			templateVo.TemplateKindWhatsappAdmin.Equals(otherValue) ||
// 			templateVo.TemplateKindWhatsappRequest.Equals(otherValue) {
// 			return regexp.MustCompile("^[a-zA-Z0-9_:\\-]+$").
// 				MatchString(field.String())
// 		}

// 		return true
// 	}
// 	return makeValidValidation(v, fn, "tmpl_name")
// }

// func registerTemplateCodeIdValidation(v *Validator) (err error) {
// 	fn := func(fl govalidator.FieldLevel) bool {
// 		return regexp.MustCompile("^[a-zA-Z0-9_:\\-]+$").
// 			MatchString(fl.Field().String())
// 	}
// 	return makeValidValidation(v, fn, "tmpl_code")
// }
