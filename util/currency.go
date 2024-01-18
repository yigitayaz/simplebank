package util

var SupportedCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"TRY": true,
}

func IsSupportedCurrency(currency string) bool {
	_, ok := SupportedCurrencies[currency]
	return ok
}
