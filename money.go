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
func (m *Money) Equals(om *Money) (bool, error) {
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

// IsZero returns true if the Money value is equals to zero,
// false otherwise
func (m *Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive returns true if the Money value is positive,
// false otherwise
func (m *Money) IsPositive() bool {
	return m.amount > 0
}

// IsNegative returns true if the Money value is negative,
// false otherwise
func (m *Money) IsNegative() bool {
	return m.amount < 0
}

// Absolute returns a new Money from a given Money
// using the absolute monetary value
func (m *Money) Absolute() *Money {
	panic("not implemented")
}

// Negative returns a new Money from a given Money
// using the negative monetary value
func (m *Money) Negative() *Money {
	panic("not implemented")
}

// Add returns a new Money with the value representing the sum of Self and Other Money
// using the negative monetary value,
// or returns a money.ErrCurrencyMismatch if currencies don't match
func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}
	panic("not implemented")
}

// Subtract returns a new Money with the value representing the difference of Self and Other Money
// using the negative monetary value,
// or returns a money.ErrCurrencyMismatch if currencies don't match
func (m *Money) Subtract(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}
	panic("not implemented")
}

// Multiply returns a new Money with the value representing Self value multiplied with multiplier
func (m *Money) Multiply(mul int64) *Money {
	panic("not implemented")
}

// Round returns a new Money with the value rounded
func (m *Money) Round(mul int64) *Money {
	panic("not implemented")
}

// Split returns a new a slice of Monies with the Self value split in given number.
// The leftover after the division will be distributed round-robin amongst the parties.
// Parties listed first will likely receive more cents than ones listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	panic("not implemented")
}

// SplitWithReminder returns a new a slice of Monies with the Self value split equally in given number.
// The reminder is returned as separate non nil Money giving the handling to the caller.
func (m *Money) SplitWithReminder(n int) ([]*Money, *Money, error) {
	panic("not implemented")
}

// Allocate returns a slice of Monies with the Self value split in given rations.
// After allocation the reminder is distributed equally amongst the parties with round-robin principle.
func (m *Money) Allocate(rs ...int) ([]*Money, error) {
	panic("not implemented")
}

// AllocateWithReminder returns a slice of Monies with the Self value split in given rations.
// After allocation the reminder is returned as a separate non nil Money giving the handling to the caller.
func (m *Money) AllocateWithReminder(rs ...int) ([]*Money, *Money, error) {
	panic("not implemented")
}

// Display displays the Money as a string in given Currency
func (m *Money) Display() string {
	panic("not implemented")
}

// UnmarshalJSON is implementation of json.Unmarshaller
func (m *Money) UnmarshalJSON(b []byte) error {
	return UnmarshalJSON(m, b)
}

// MarshalJSON is implementation of json.Unmarshaller
func (m Money) MarshalJSON() ([]byte, error) {
	return MarshalJSON(m)
}

// ...
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
