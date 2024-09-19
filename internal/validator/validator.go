package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	govalidator "github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	ierrors "github.com/higansama/xyz-multi-finance/internal/errors"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Validator struct {
	engine     *govalidator.Validate
	translator ut.Translator
}

type genericValidatorData struct {
	GenericField string
}

func ConfigValidator() *Validator {
	engine := binding.Validator.Engine().(*govalidator.Validate)

	// re-format StructFields name, so it can be customized later easily
	engine.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := utils.GetFieldNameFromStructField(fld)

		return "::" + name + "||" + fld.Name + "::"
	})

	english := en.New()
	uni := ut.New(english, english)
	translator, _ := uni.GetTranslator("en")
	transErr := entranslations.RegisterDefaultTranslations(engine, translator)

	if transErr != nil {
		errMsg := fmt.Sprintf("RegisterDefaultTranslations Error: %v", transErr)
		defer log.Fatal().Msg(errMsg)
		panic(errMsg)
	}

	v := &Validator{
		engine:     engine,
		translator: translator,
	}

	RegisterCustomValidator(v)
	// RegisterSpecificValidator(v)
	RegisterValidatorTranslations(v)

	return v
}

func (r *Validator) Validate(obj any) error {
	if !utils.IsStructAlike(obj) {
		errMsg := fmt.Sprintf(
			"validator@Validate (Validator value not supported, because %v is not struct type)",
			reflect.TypeOf(obj))
		defer log.Fatal().Msg(errMsg)
		panic(errMsg)
	}

	err := r.engine.Struct(obj)

	return FormatValidationError(r, obj, err)
}

func (r *Validator) ValidateSingleField(value string, tag string, field string) error {
	rules := map[string]string{
		"GenericField": tag,
	}

	data := genericValidatorData{
		GenericField: value,
	}

	r.engine.RegisterStructValidationMapRules(rules, genericValidatorData{})

	err := r.engine.Struct(data)

	if err != nil {
		err := FormatValidationError(r, data, err)
		var vErrs ierrors.ValidationError
		if !errors.As(err, &vErrs) {
			return err
		}

		fErr := vErrs.Errors[0]
		vErrs.Errors = nil

		title := cases.Title(language.English)
		fErr.Msg = strings.Replace(fErr.Msg, "GenericField", title.String(field), 1)
		vErrs.Errors = append(vErrs.Errors, fErr)

		return vErrs
	}

	return nil
}

func FormatValidationError(v *Validator, obj any, err error) error {
	res := ierrors.ValidationError{}
	res.Errors = make([]ierrors.FieldError, 0)

	var verrs govalidator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}

	title := cases.Title(language.English)

	for _, fe := range verrs {
		ns := fe.StructNamespace()
		nsk := strings.Split(ns, ".")[1:] // remove struct name
		key := strings.Join(nsk, ".")
		ins := len(nsk) > 1 // in nested struct
		if len(nsk) > 1 {
			nsk = nsk[:len(nsk)-1] // except last one
		}
		baseKey := strings.Join(nsk, ".")

		// @TODO: FIX ASAP, error Field should be actual field, don't get from `attr`
		fld := utils.GetFieldNameForNamespace(obj, key, "")
		fnn := strings.Split(fld, ".")
		fnnd := fld
		if len(fnn) > 1 {
			fnnd = fnn[len(fnn)-1]
		}

		errResult := ierrors.FieldError{
			Field: fld,
			Msg:   fe.Translate(v.translator),
			Tag:   fe.Tag(),
		}

		re := regexp.MustCompile(`::(?P<field>[0-9a-zA-Z_ ]+)\|\|(?P<tag>[0-9a-zA-Z_ ]+)?::(\[(?P<index>[0-9]+)\])?`)
		result := re.FindAllStringSubmatch(errResult.Msg, 99)

		for _, matches := range result {
			rpc := matches[1]
			frpc := rpc
			if matches[4] != "" {
				frpc = matches[4]
			}

			if frpc != fnnd {
				bk := baseKey + "."
				if !ins {
					bk = ""
				}
				// get field name for provided param
				mk := utils.GetFieldNameForNamespace(obj, bk+rpc, "")
				nss := strings.Split(mk, ".")
				rpc = mk
				if len(nss) > 1 {
					rpc = nss[len(nss)-1]
				}
			}

			rpc = title.String(strings.ReplaceAll(rpc, "_", " "))
			errResult.Msg = strings.Replace(errResult.Msg, matches[0], rpc, 1)
		}

		res.Errors = append(res.Errors, errResult)
	}

	if len(res.Errors) == 0 {
		return nil
	}

	return res
}

func makeValidValidation(v *Validator, fn func(fl govalidator.FieldLevel) bool, vn string) (err error) {
	err = v.engine.RegisterValidation(vn, fn, true)
	if err != nil {
		return
	}

	err = v.engine.RegisterTranslation(vn, v.translator, func(ut ut.Translator) error {
		return ut.Add(vn, "{0} is invalid", false)
	}, func(ut ut.Translator, fe govalidator.FieldError) string {
		t, _ := ut.T(vn, fe.Field())
		return t
	})
	if err != nil {
		return
	}

	return
}

var (
	oneofValsCache       = map[string][]string{}
	oneofValsCacheRWLock = sync.RWMutex{}
)

var splitParamsRegexString = `'[^']*'|\S+`
var splitParamsRegex = regexp.MustCompile(splitParamsRegexString)

func parseOneOfParam2(s string) []string {
	oneofValsCacheRWLock.RLock()
	vals, ok := oneofValsCache[s]
	oneofValsCacheRWLock.RUnlock()
	if !ok {
		oneofValsCacheRWLock.Lock()
		vals = splitParamsRegex.FindAllString(s, -1)
		for i := 0; i < len(vals); i++ {
			vals[i] = strings.Replace(vals[i], "'", "", -1)
		}
		oneofValsCache[s] = vals
		oneofValsCacheRWLock.Unlock()
	}
	return vals
}

func requireCheckFieldValue(
	fl govalidator.FieldLevel,
	param string,
	value string,
	defaultNotFoundValue bool,
) bool {
	field, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
	if !found {
		return defaultNotFoundValue
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() == asInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return field.Uint() == asUint(value)
	case reflect.Float32:
		return field.Float() == asFloat32(value)
	case reflect.Float64:
		return field.Float() == asFloat64(value)
	case reflect.Slice, reflect.Map, reflect.Array:
		return int64(field.Len()) == asInt(value)
	case reflect.Bool:
		return field.Bool() == asBool(value)
	}

	return field.String() == value
}

// hasValue is the validation function for validating if the current field's value is not the default static value.
func hasValue(fl govalidator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		//if fl.(*validate).fldIsPointer && field.Interface() != nil {
		//	return true
		//}
		return field.IsValid() && !field.IsZero()
	}
}
