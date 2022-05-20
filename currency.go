package money

import "strings"

// Currency is an ISO 4217 currency designator
type Currency struct {
	code string
}

// GetCurrencyForCode returns the currency for a given code
// or an money.ErrUnsupportedCurrency error if currency could not be found.
func GetCurrencyForCode(code string) (*Currency, error) {
	currency, prs := currencies[code]
	if !prs {
		return nil, ErrUnsupportedCurrency
	}
	return currency, nil
}

// MustGetCurrencyForCode returns the currency for a given code
// or panics if currency could not be found
func MustGetCurrencyForCode(code string) *Currency {
	currency, err := GetCurrencyForCode(code)
	if err != nil {
		panic(err)
	}
	return currency
}

// Code returns the currency's code
func (c *Currency) Code() string {
	return c.code
}

// Equals returns true, if both Codes are equal,
// false otherwise
func (c *Currency) Equals(oc Currency) bool {
	return c.code == oc.code
}

// String returns the string representation of the currency
func (c *Currency) String() string {
	return c.code
}

func newCurrency(code string) *Currency {
	return &Currency{code: strings.ToUpper(code)}
}

// Currencies holds all the supported currencies for a given code
type Currencies map[string]*Currency

var currencies = Currencies{
	EUR: newCurrency("EUR"),
	USD: newCurrency("USD"),
}
