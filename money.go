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

// SameCurrency returns true, if given Money is equals by currency,
// false otherwise
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

// Equals returns true, if given Money is equals by currency and amount,
// false with error money.ErrCurrencyMismatch if currencies don't match,
// or false if currency matches, but the amount differs
func (m *Money) Equal(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}
	return m.compare(om) == 0, nil
}

// GreaterThan checks wether Money value is greater than the others,
// returns an money.ErrCurrencyMismatch error if the currencies don't match
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}
	return m.compare(om) == 1, nil
}

// GreaterThanOrEqual checks wether Money value is greater than or equal the others,
// returns an money.ErrCurrencyMismatch error if the currencies don't match
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}
	return m.compare(om) >= 0, nil
}

// LessThan checks wether Money value is less than the others,
// returns an money.ErrCurrencyMismatch error if the currencies don't match
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}
	return m.compare(om) == -1, nil
}

// LessThanOrEqual checks wether Money value is less than or equal the others,
// returns an money.ErrCurrencyMismatch error if the currencies don't match
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}
	return m.compare(om) <= 0, nil
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return ErrCurrencyMismatch
	}
	return nil
}

func (m *Money) compare(om *Money) int {
	switch {
	case m.amount > om.amount:
		return 1
	case m.amount < om.amount:
		return -1
	}
	return 0
}
