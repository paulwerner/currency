package money

type Money struct {
	amount   int64
	currency *Currency
}

// New creates a new Money value with the given amount and code
func New(amount int64, code string) *Money {
	return &Money{
		amount:   amount,
		currency: getOrDefault(code),
	}
}

// Currency returns the currency for the monetary value
func (m *Money) Currency() *Currency {
	return m.currency
}

// CurrencyCode returns the currency code for the monetary value
func (m *Money) CurrencyCode() string {
	return m.currency.code
}

// Amount returns the amount in the fractional monetary unit
func (m *Money) Amount() int64 {
	return m.amount
}
