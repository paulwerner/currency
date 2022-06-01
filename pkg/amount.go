package currency

import "errors"

var (
	ErrCurrencyMismatch = errors.New("money: currency mismatch")
	ErrInvalidOperation = errors.New("money: invalid operation")
	ErrSplitNegative    = errors.New("money: split must be positive")
	ErrNoRatioSpecified = errors.New("money: no ratio specified")
)

type Amount struct {
	value    int
	currency *Currency
}

func NewAmount(v int, cur Currency) (*Amount, error) {
	return &Amount{
		value:    v,
		currency: &cur,
	}, nil
}

func NewFromISO(v int, iso string) (*Amount, error) {
	cur, err := CurrencyFromISO(iso)
	if err != nil {
		return nil, err
	}
	return &Amount{
		value:    v,
		currency: cur,
	}, nil
}

func (a *Amount) Currency() *Currency {
	return a.currency
}

func (a *Amount) Amount() int {
	return a.value
}

func (a *Amount) SameCurrency(om *Amount) bool {
	return a.currency.Equals(om.currency)
}

func (a *Amount) Add(om *Amount) (*Amount, error) {
	if err := a.assertSameCurrency(om); err != nil {
		return nil, err
	}

	z, ok := add(a.value, om.value)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{
		value:    z,
		currency: a.currency,
	}, nil
}

func (a *Amount) Sub(om *Amount) (*Amount, error) {
	if err := a.assertSameCurrency(om); err != nil {
		return nil, err
	}

	z, ok := sub(a.value, om.value)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{
		value:    z,
		currency: a.currency,
	}, nil
}

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

func (a *Amount) Split(n int) ([]*Amount, *Amount, error) {
	if n <= 0 {
		return nil, nil, ErrSplitNegative
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

func (a *Amount) Alloc(rs ...int) ([]*Amount, *Amount, error) {
	if len(rs) == 0 {
		return nil, nil, ErrNoRatioSpecified
	}

	// sum of ratios
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total int
	var ms []*Amount

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

func (a *Amount) Round(k Kind) (*Amount, error) {
	s, i := k.Rounding(a.currency)
	r, ok := round(a.value, s, i)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Amount{value: r, currency: a.currency}, nil
}

func (a *Amount) Display() string {
	panic("not implemented")
}

func (a *Amount) String() string {
	panic("not implemented")
}

func (a *Amount) MarshalJSON() ([]byte, error) {
	panic("not implemented")
}
func (a *Amount) UnmarshalJSON([]byte) error {
	panic("not implemented")
}

func (a *Amount) Equals(oa *Amount) bool {
	return a.value == oa.value &&
		a.currency.Equals(oa.currency)
}

func (a *Amount) GreaterThan(oa *Amount) (bool, error) {
	panic("not implemented")
}

func (a *Amount) GreaterThanOrEqual(oa *Amount) (bool, error) {
	panic("not implemented")
}

func (a *Amount) LessThan(oa *Amount) (bool, error) {
	panic("not implemented")
}

func (a *Amount) LessThanOrEqual(oa *Amount) (bool, error) {
	panic("not implemented")
}

func (a *Amount) IsPositive() bool {
	return a.value >= 0
}

func (a *Amount) IsZero() bool {
	return a.value == 0
}

func (a *Amount) IsNegative() bool {
	return a.value < 0
}

func (a *Amount) Abs() (bool, error) {
	panic("not implemented")
}

func (a *Amount) Neg() (bool, error) {
	panic("not implemented")
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
