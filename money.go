package money

// Money holds the monetary unit and currency
type Money struct {
	amount   int64
	currency *Currency
}

// New creates a new Money value for the given amount
// and currency code
func New(amount int64, code string) *Money {
	return &Money{
		amount:   amount,
		currency: &Currency{code},
	}
}

// Currency returns the currency
func (m *Money) Currency() *Currency {
	return m.currency
}

// CurrencyCode returns a copy of the currency's code
func (m *Money) CurrencyCode() string {
	return m.currency.code
}

// Amount returns a copy of the internal monetary value as Amount
func (m *Money) Amount() int64 {
	return m.amount
}
