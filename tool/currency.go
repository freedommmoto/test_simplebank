package tool

const (
	USDcurrency = "USD"
	THBcurrency = "THB"
)

func IsSupportedCurrencyOrNot(currency string) bool {
	switch currency {
	case USDcurrency, THBcurrency:
		return true
	}
	return false
}
