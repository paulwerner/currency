package money

import "strings"

type Currency struct {
	code     string
	fraction int
}

// equals returns true, if both currencies have the same code,
// false otherwise
func (c *Currency) equals(oc *Currency) bool {
	return c.code == oc.code
}

// GetCurrency returns a Currency for a given code if supported,
// or ErrUnsupportedCurrency otherwise
func GetCurrency(code string) (*Currency, error) {
	return supported.currencyByCode(code)
}

// MustGetCurrency returns a Currency for a given code if supported,
// or panics otherwise
func MustGetCurrency(code string) *Currency {
	c, err := supported.currencyByCode(code)
	if err != nil {
		panic(err)
	}
	return c
}

type currencies map[string]*Currency

func (c currencies) currencyByCode(code string) (*Currency, error) {
	v, ok := c[code]
	if !ok {
		return nil, ErrUnsupportedCurrency
	}
	return v, nil
}

func getOrDefault(code string) *Currency {
	v, ok := supported[code]
	if ok {
		return v
	}
	return supported[USD]
}

func newCurrency(code string, fraction int) *Currency {
	return &Currency{
		code:     strings.ToUpper(code),
		fraction: fraction,
	}
}

// the supported currencies
var supported = currencies{
	USD: newCurrency("USD", 2),
	EUR: newCurrency("EUR", 2),
}
