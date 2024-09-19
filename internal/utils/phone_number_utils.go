package utils

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"strings"
)

func ConstructPhoneNumber(countryCode string, defCode string, number string) (string, string) {
	number = ReformatPhoneNumber(number)
	if countryCode == "" {
		countryCode = GetCountryCodeFromPhoneNumber(number, defCode)
	}

	return countryCode, ReformatPhoneNumber(countryCode + strings.TrimPrefix(number, countryCode))
}

func ReformatPhoneNumber(number string) string {
	hasPlusLookup := strings.HasPrefix(number, "+")
	number = phonenumbers.NormalizeDigitsOnly(number)

	knownReplacer := map[string][]string{
		"62": {"620", "62"},
	}
	mapped := false
	for k, v := range knownReplacer {
		for _, px := range v {
			if strings.HasPrefix(number, px) {
				mapped = true
				number = k + strings.TrimPrefix(number, px)
				break
			}
		}
		if mapped {
			break
		}
	}
	if hasPlusLookup || mapped {
		number = "+" + number
	}

	return number
}

func GetCountryCodeFromPhoneNumber(number string, def string) string {
	num, err := phonenumbers.Parse(number, "")
	if err != nil {
		return def
	}

	return fmt.Sprintf("+%d", *num.CountryCode)
}
