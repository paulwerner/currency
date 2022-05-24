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

func (a *Amount) Add(oa *Amount) (*Amount, error) {
	panic("not implemented")
}

func (a *Amount) Sub(oa *Amount) (*Amount, error) {
	panic("not implemented")
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
