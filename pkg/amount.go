package money

type Amount struct {
	val      int64
	currency Currency
}

// NewAmount creates a new amount for the given val and ISO 4217 currency code.
// If returns an error if the ISO code is not supported.
func NewAmount(val int64, isoCode string) (*Amount, error) {
	cur, err := ParseISO(isoCode)
	if err != nil {
		return nil, err
	}
	return &Amount{val: val, currency: cur}, nil
}

// Currency reports the amount's currency
func (m *Amount) Currency() Currency {
	return m.currency
}

func (m *Amount) Add(om *Amount) (*Amount, error) {
	panic("not implemented")
}

func (m *Amount) Sub(om *Amount) (*Amount, error) {
	panic("not implemented")
}

func (m *Amount) Div(d int64) (*Amount, error) {
	panic("not implemented")
}

func (m *Amount) Mul(mul int64) (*Amount, error) {
	panic("not implemented")
}

func (m *Amount) Round() (*Amount, error) {
	panic("not implemented")
}

func (m *Amount) Split(p int64) ([]*Amount, error) {
	panic("not implemented")
}

func (m *Amount) Alloc(ps ...int64) ([]*Amount, error) {
	panic("not implemented")
}

func (m *Amount) Display() string {
	panic("not implemented")
}