package money

import "strings"

type Currency struct {
	code string
}

func GetCurrencyForCode(code string) (*Currency, error) {
	currency, prs := currencies[code]
	if !prs {
		return nil, ErrUnsupportedCurrency
	}
	return currency, nil
}

func MustGetCurrencyForCode(code string) *Currency {
	currency, err := GetCurrencyForCode(code)
	if err != nil {
		panic(err)
	}
	return currency
}

func (c *Currency) Code() string {
	return c.code
}

func (c *Currency) Equals(oc Currency) bool {
	return c.code == oc.code
}

func (c *Currency) String() string {
	return c.code
}

func newCurrency(code string) *Currency {
	return &Currency{code: strings.ToUpper(code)}
}

type Currencies map[string]*Currency

var currencies = Currencies{
	EUR: newCurrency("EUR"),
	USD: newCurrency("USD"),
}
