package currency

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrCurrencyMismatch   = errors.New("currency: currencies do not match")
	ErrInvalidOperation   = errors.New("currency: invalid operation")
	ErrInvalidSplitNumber = errors.New("currency: split number must be positive")
	ErrNoRatiosSpecified  = errors.New("currency: no allocation ratios specified")
	ErrInvalidJSON        = errors.New("currency: invalid json")
)

// Amount is a monetary value: an integer amount in the currency's minor
// units (e.g. cents for EUR, yen for JPY) paired with a Currency.
//
// The zero value is 0 in the undefined currency XXX and is ready to use.
// All arithmetic is overflow-checked and returns an error instead of
// silently wrapping around. Amounts are immutable; operations return a
// new Amount
type Amount struct {
	value    int
	currency Currency
}

// NewAmount returns an Amount of v minor units of the given currency
func NewAmount(v int, cur Currency) *Amount {
	return &Amount{
		value:    v,
		currency: cur,
	}
}

// NewFromISO returns an Amount of v minor units of the currency identified
// by the given 3-letter ISO 4217 code
func NewFromISO(v int, iso string) (*Amount, error) {
	cur, err := CurrencyFromISO(iso)
	if err != nil {
		return nil, err
	}
	return &Amount{
		value:    v,
		currency: *cur,
	}, nil
}

// Currency reports the amount's currency
func (a *Amount) Currency() Currency {
	return a.currency
}

// Amount reports the amount's value in the currency's minor units
func (a *Amount) Amount() int {
	return a.value
}

// SameCurrency reports whether both amounts are in the same currency
func (a *Amount) SameCurrency(oa *Amount) bool {
	return a.currency.Equals(oa.currency)
}

// Add returns the sum a + oa. It returns ErrCurrencyMismatch if the
// currencies differ and ErrInvalidOperation on overflow
func (a *Amount) Add(oa *Amount) (*Amount, error) {
	if err := a.assertSameCurrency(oa); err != nil {
		return nil, err
	}

	z, ok := add(a.value, oa.value)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{
		value:    z,
		currency: a.currency,
	}, nil
}

// Sub returns the difference a - oa. It returns ErrCurrencyMismatch if the
// currencies differ and ErrInvalidOperation on overflow
func (a *Amount) Sub(oa *Amount) (*Amount, error) {
	if err := a.assertSameCurrency(oa); err != nil {
		return nil, err
	}

	z, ok := sub(a.value, oa.value)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{
		value:    z,
		currency: a.currency,
	}, nil
}

// Mul returns the product a * n. It returns ErrInvalidOperation on overflow
func (a *Amount) Mul(n int) (*Amount, error) {
	z, ok := mul(a.value, n)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{
		value:    z,
		currency: a.currency,
	}, nil
}

// Split divides the amount into n equal parts and a remainder holding the
// indivisible rest, such that the parts and the remainder sum up to the
// original amount. It returns ErrInvalidSplitNumber if n is not positive
func (a *Amount) Split(n int) ([]*Amount, *Amount, error) {
	if n <= 0 {
		return nil, nil, ErrInvalidSplitNumber
	}

	z, ok := div(a.value, n)
	if !ok {
		return nil, nil, ErrInvalidOperation
	}

	ms := make([]*Amount, n)
	for i := 0; i < n; i++ {
		ms[i] = &Amount{value: z, currency: a.currency}
	}
	r, ok := mod(a.value, n)
	if !ok {
		return nil, nil, ErrInvalidOperation
	}
	return ms, &Amount{value: r, currency: a.currency}, nil
}

// Alloc distributes the amount according to the given ratios, e.g. 50, 30,
// 20, and returns one part per ratio plus a leftover holding the
// indivisible rest, such that parts and leftover sum up to the original
// amount. Ratios must not be negative and must not sum up to zero
func (a *Amount) Alloc(rs ...int) ([]*Amount, *Amount, error) {
	if len(rs) == 0 {
		return nil, nil, ErrNoRatiosSpecified
	}

	// sum of ratios
	var sum int
	var ok bool
	for _, r := range rs {
		sum, ok = add(sum, r)
		if !ok {
			return nil, nil, ErrInvalidOperation
		}
	}

	var total int
	ms := make([]*Amount, 0, len(rs))

	for _, r := range rs {
		alloc, ok := alloc(a.value, r, sum)
		if !ok {
			return nil, nil, ErrInvalidOperation
		}
		party := &Amount{
			value:    alloc,
			currency: a.currency,
		}
		ms = append(ms, party)
		total, ok = add(total, alloc)
		if !ok {
			return nil, nil, ErrInvalidOperation
		}
	}

	// leftover
	lo, ok := sub(a.value, total)
	if !ok {
		return nil, nil, ErrInvalidOperation
	}
	return ms, &Amount{value: lo, currency: a.currency}, nil
}

// Round rounds the amount according to the rounding rules of the given
// kind, with ties rounded away from zero. For Standard and Accounting this
// rounds to the currency's smallest unit (a no-op for amounts already in
// minor units); for Cash it rounds to the currency's cash rounding
// increment, e.g. to 0.05 for CHF
func (a *Amount) Round(k Kind) (*Amount, error) {
	step, ok := a.roundingStep(k)
	if !ok {
		return nil, ErrInvalidOperation
	}
	r, ok := round(a.value, step)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{value: r, currency: a.currency}, nil
}

// roundingStep reports the kind's rounding increment in minor units
func (a *Amount) roundingStep(k Kind) (int, bool) {
	stdScale, _ := Standard.Rounding(a.currency)
	scale, inc := k.Rounding(a.currency)
	if scale > stdScale {
		// rounding scales beyond the currency's minor unit cannot be
		// represented by an integer amount of minor units
		scale = stdScale
	}
	f, ok := pow(10, stdScale-scale)
	if !ok {
		return 0, false
	}
	return mul(inc, f)
}

// String returns a locale-independent representation of the amount in the
// form "EUR 12.34", using the currency's standard scale
func (a *Amount) String() string {
	return a.currency.String() + " " + a.formatValue()
}

// formatValue renders the plain decimal value using the standard scale,
// e.g. 1234 EUR -> "12.34", 1234 JPY -> "1234"
func (a *Amount) formatValue() string {
	scale, _ := Standard.Rounding(a.currency)
	if scale == 0 {
		// also avoids negating loBound below, which would overflow
		return strconv.Itoa(a.value)
	}
	exp, ok := pow(10, scale)
	if !ok {
		return strconv.Itoa(a.value)
	}

	sign := ""
	// exp >= 10 here, so |whole| and |frac| are far from the integer
	// bounds and safe to negate
	whole, frac := a.value/exp, a.value%exp
	if a.value < 0 {
		sign = "-"
		whole, frac = -whole, -frac
	}
	return fmt.Sprintf("%s%d.%0*d", sign, whole, scale, frac)
}

// amountJSON is the wire representation of an Amount
type amountJSON struct {
	Amount   *int64 `json:"amount"`
	Currency string `json:"currency"`
}

// MarshalJSON encodes the amount as {"amount":<minor units>,"currency":"<ISO code>"}
func (a *Amount) MarshalJSON() ([]byte, error) {
	return json.Marshal(amountJSON{
		Amount:   ptr(int64(a.value)),
		Currency: a.currency.Code(),
	})
}

// UnmarshalJSON decodes an amount from {"amount":<minor units>,"currency":"<ISO code>"}.
// Both fields are required
func (a *Amount) UnmarshalJSON(b []byte) error {
	var aj amountJSON
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&aj); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if aj.Amount == nil {
		return fmt.Errorf("%w: missing field \"amount\"", ErrInvalidJSON)
	}
	if aj.Currency == "" {
		return fmt.Errorf("%w: missing field \"currency\"", ErrInvalidJSON)
	}
	v := int(*aj.Amount)
	if int64(v) != *aj.Amount {
		return fmt.Errorf("%w: amount out of range", ErrInvalidJSON)
	}
	cur, err := CurrencyFromISO(aj.Currency)
	if err != nil {
		return err
	}
	a.value = v
	a.currency = *cur
	return nil
}

func ptr[T any](v T) *T { return &v }

// Equals reports whether both amounts have the same value and currency
func (a *Amount) Equals(oa *Amount) bool {
	return a.value == oa.value &&
		a.currency.Equals(oa.currency)
}

// Cmp compares a and oa and returns -1 if a is less than oa, 0 if they are
// equal, and +1 if a is greater than oa. It returns ErrCurrencyMismatch if
// the currencies differ
func (a *Amount) Cmp(oa *Amount) (int, error) {
	if err := a.assertSameCurrency(oa); err != nil {
		return 0, err
	}
	switch {
	case a.value < oa.value:
		return -1, nil
	case a.value > oa.value:
		return 1, nil
	default:
		return 0, nil
	}
}

// GreaterThan reports whether a is greater than oa.
// It returns ErrCurrencyMismatch if the currencies differ
func (a *Amount) GreaterThan(oa *Amount) (bool, error) {
	c, err := a.Cmp(oa)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

// GreaterThanOrEqual reports whether a is greater than or equal to oa.
// It returns ErrCurrencyMismatch if the currencies differ
func (a *Amount) GreaterThanOrEqual(oa *Amount) (bool, error) {
	c, err := a.Cmp(oa)
	if err != nil {
		return false, err
	}
	return c >= 0, nil
}

// LessThan reports whether a is less than oa.
// It returns ErrCurrencyMismatch if the currencies differ
func (a *Amount) LessThan(oa *Amount) (bool, error) {
	c, err := a.Cmp(oa)
	if err != nil {
		return false, err
	}
	return c < 0, nil
}

// LessThanOrEqual reports whether a is less than or equal to oa.
// It returns ErrCurrencyMismatch if the currencies differ
func (a *Amount) LessThanOrEqual(oa *Amount) (bool, error) {
	c, err := a.Cmp(oa)
	if err != nil {
		return false, err
	}
	return c <= 0, nil
}

// IsPositive reports whether the amount is greater than zero
func (a *Amount) IsPositive() bool {
	return a.value > 0
}

// IsZero reports whether the amount is zero
func (a *Amount) IsZero() bool {
	return a.value == 0
}

// IsNegative reports whether the amount is less than zero
func (a *Amount) IsNegative() bool {
	return a.value < 0
}

// Abs returns the absolute value of the amount.
// It returns ErrInvalidOperation on overflow
func (a *Amount) Abs() (*Amount, error) {
	v, ok := abs(a.value)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{value: v, currency: a.currency}, nil
}

// Neg returns the amount with its sign inverted.
// It returns ErrInvalidOperation on overflow
func (a *Amount) Neg() (*Amount, error) {
	v, ok := neg(a.value)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{value: v, currency: a.currency}, nil
}

//
// PRIVATE
//

func (a *Amount) assertSameCurrency(oa *Amount) error {
	if !a.currency.Equals(oa.currency) {
		return ErrCurrencyMismatch
	}
	return nil
}
