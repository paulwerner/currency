package money

type Amount = int64

type Money struct {
	amount   Amount
	currency *Currency
}

// New creates a new Money value with the given amount and code
func New(amount Amount, code string) *Money {
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
func (m *Money) Amount() Amount {
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

// Abs returns a new Money from a given Money
// using the absolute monetary value
func (m *Money) Abs() *Money {
	return &Money{amount: calc.abs(m.amount), currency: m.currency}
}

// Neg returns a new Money from a given Money
// using the negative monetary value
func (m *Money) Neg() *Money {
	return &Money{amount: calc.neg(m.amount), currency: m.currency}
}

// Add returns a new Money with the value representing the sum of Self and Other Money
// using the negative monetary value,
// or returns a money.ErrCurrencyMismatch if currencies don't match
func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}
	return &Money{amount: calc.add(m.amount, om.amount), currency: m.currency}, nil
}

// Sub returns a new Money with the value representing the difference of Self and Other Money
// using the negative monetary value,
// or returns a money.ErrCurrencyMismatch if currencies don't match
func (m *Money) Sub(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}
	return &Money{amount: calc.sub(m.amount, om.amount), currency: m.currency}, nil
}

// Mul returns a new Money with the value representing Self value multiplied with multiplier
func (m *Money) Mul(mul int64) *Money {
	return &Money{amount: calc.mul(m.amount, mul), currency: m.currency}
}

// Round returns a new Money with the value rounded
func (m *Money) Round() *Money {
	return &Money{amount: calc.round(m.amount, m.currency.fraction), currency: m.currency}
}

// Split returns a new a slice of Monies with the Self value split in given number.
// The leftover after the division will be distributed round-robin amongst the parties.
// Parties listed first will likely receive more cents than ones listed later.
func (m *Money) Split(n int) ([]*Money, error) {
	ms, l, err := m.split(n)
	if err != nil {
		return nil, err
	}
	// distribute leftover amongst first parties
	v := Amount(1)
	if m.amount < 0 {
		v = -1
	}
	for p := 0; l != 0; p++ {
		ms[p].amount = calc.add(ms[p].amount, v)
		l--
	}
	return ms, nil
}

// SplitWithReminder returns a new a slice of Monies with the Self value split equally in given number.
// The reminder is returned as separate non nil Money giving the handling to the caller.
func (m *Money) SplitWithReminder(n int) ([]*Money, *Money, error) {
	ms, l, err := m.split(n)
	if err != nil {
		return nil, nil, err
	}
	return ms, &Money{amount: l, currency: m.currency}, nil
}

// Alloc returns a slice of Monies with the Self value split in given rations.
// After allocation the reminder is distributed equally amongst the parties with round-robin principle.
func (m *Money) Alloc(rs ...int) ([]*Money, error) {
	ms, lo, err := m.alloc(rs...)
	if err != nil {
		return nil, err
	}

	// distribute leftover equally amongst first parties
	sub := Amount(1)
	if lo < 0 {
		sub = -sub
	}

	for p := 0; lo != 0; p++ {
		ms[p].amount = calc.add(ms[p].amount, sub)
		lo -= sub
	}

	return ms, nil
}

// AllocWithReminder returns a slice of Monies with the Self value split in given rations.
// After allocation the reminder is returned as a separate non nil Money giving the handling to the caller.
func (m *Money) AllocWithReminder(rs ...int) ([]*Money, *Money, error) {
	ms, lo, err := m.alloc(rs...)
	if err != nil {
		return nil, nil, err
	}
	return ms, &Money{amount: lo, currency: m.currency}, nil
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

func (m *Money) split(n int) ([]*Money, Amount, error) {
	if n <= 0 {
		return nil, 0, ErrSplitNotPositive
	}

	a := calc.div(m.amount, Amount(n))
	ms := make([]*Money, n)

	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: a, currency: m.currency}
	}

	r := calc.mod(m.amount, Amount(n))
	l := calc.abs(r)

	return ms, l, nil
}

func (m *Money) alloc(rs ...int) ([]*Money, Amount, error) {
	if len(rs) == 0 {
		return nil, 0, ErrNoRatiosSpecified
	}

	var sum int
	for _, r := range rs {
		sum += r
	}

	var total Amount
	ms := make([]*Money, 0, len(rs))
	for _, r := range rs {
		p := &Money{
			amount:   calc.alloc(m.amount, r, sum),
			currency: m.currency,
		}
		ms = append(ms, p)
		total += p.amount
	}

	lo := m.amount - total

	return ms, lo, nil
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
