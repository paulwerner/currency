package money

import "errors"

// constants for  codes according to ISO 4217
const (
	EUR = "EUR"
)

// errors
var (
	ErrUnsupportedCurrency = errors.New("unsupported currency")
)
