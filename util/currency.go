package util

const (
	USD = "USD"
	EUR = "EUR"
)

// IsSupportCurrency returns true if the currency is supported.
func IsSupportCurrency(currency string) bool {
	switch currency {
	case USD, EUR:
		return true
	}
	return false
}
