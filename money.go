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

