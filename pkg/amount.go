package money

import "errors"

type Value = int64
type Amount struct {
	val      Value
	currency Currency
}

// NewAmount creates a new amount for the given val and ISO 4217 currency code.
// It returns an error if the ISO code is not well-formed or not recognized.
func NewAmount(val Value, isoCode string) (*Amount, error) {
	cur, err := ParseISO(isoCode)
	if err != nil {
		return nil, err
	}
	return &Amount{val: val, currency: cur}, nil
}

// Value reports the value of the current amount
func (a *Amount) Value() Value {
	return a.val
}

// Currency reports the amount's currency
func (a *Amount) Currency() Currency {
	return a.currency
}

// Equals returns true, if the amount value and currency equals,
// false otherwise
func (a *Amount) Equals(oa *Amount) bool {
	return a.val == oa.val && a.currency.Equals(&oa.currency)
}

// Add creates a new Amount representing the sum with the given amount, 
// or returns an error if the currencies mismatch. 
func (a *Amount) Add(oa *Amount) (*Amount, error) {
	if err := a.assertSameCurrency(oa); err != nil {
		return nil, err
	}
	return &Amount{val: calc.add(a.val, oa.val), currency: a.currency}, nil
}

// Sub creates a new Amount representing the difference to the given amount, 
// or returns an error if the currencies mismatch.
func (a *Amount) Sub(oa *Amount) (*Amount, error) {
	if err := a.assertSameCurrency(oa); err != nil {
		return nil, err
	}
	return &Amount{val: calc.sub(a.val, oa.val), currency: a.currency}, nil
}

func (a *Amount) Div(d Value) (*Amount, error) {
	panic("not implemented")
}

func (a *Amount) Mul(aul Value) (*Amount, error) {
	panic("not implemented")
}

func (a *Amount) Round() (*Amount, error) {
	panic("not implemented")
}

func (a *Amount) Split(p Value) ([]*Amount, error) {
	panic("not implemented")
}

func (a *Amount) Alloc(ps ...Value) ([]*Amount, error) {
	panic("not implemented")
}

func (a *Amount) Display() string {
	panic("not implemented")
}

var (
	errCurrencyMismatch = errors.New("amount: currencies don't match")
)

func (a *Amount) assertSameCurrency(oa *Amount) error {
	if !a.currency.Equals(&oa.currency) {
		return errCurrencyMismatch
	}
	return nil
}
