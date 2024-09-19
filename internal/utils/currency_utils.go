package utils

import (
	"github.com/leekchan/accounting"
	"strings"
)

func FormatCurrency(value int64, curr string, precision int) string {
	defaultCurrency := "IDR"

	lc := accounting.LocaleInfo[defaultCurrency]
	if v, ok := accounting.LocaleInfo[strings.ToUpper(curr)]; ok {
		lc = v
	}

	ac := accounting.Accounting{
		Symbol:    lc.ComSymbol,
		Precision: precision,
		Thousand:  lc.ThouSep,
		Decimal:   lc.DecSep}

	res := ac.FormatMoney(value)

	// fix
	res = strings.Replace(res, "Rp.", "Rp", 1)

	return res
}
