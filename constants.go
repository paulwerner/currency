package money

import "errors"

// constants for  codes according to ISO 4217
const (
	USD = "USD"
	EUR = "EUR"
)

// errors
var (
	ErrUnsupportedCurrency = errors.New("unsupported currency")
	ErrCurrencyMismatch    = errors.New("currency don't match")
)
