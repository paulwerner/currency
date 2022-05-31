package money

import "errors"

var (
	ErrCurrencyMismatch = errors.New("money: currency mismatch")
	ErrInvalidOperation = errors.New("money: invalid operation")
	ErrSplitNegative    = errors.New("money: split must be positive")
	ErrNoRatioSpecified = errors.New("money: no ratio specified")
)

type amount = int
type Money struct {
	amount   amount
	currency *Currency
}

func New(v amount, cur Currency) (*Money, error) {
	return &Money{
		amount:   v,
		currency: &cur,
	}, nil
}

func NewFromISO(v amount, iso string) (*Money, error) {
	cur, err := CurrencyFromISO(iso)
	if err != nil {
		return nil, err
	}
	return &Money{
		amount:   v,
		currency: cur,
	}, nil
}

func (m *Money) Currency() *Currency {
	return m.currency
}

func (m *Money) Amount() amount {
	return m.amount
}

func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.Equals(om.currency)
}

func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	z, ok := add(m.amount, om.amount)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Money{
		amount:   z,
		currency: m.currency,
	}, nil
}

func (m *Money) Sub(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	z, ok := sub(m.amount, om.amount)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Money{
		amount:   z,
		currency: m.currency,
	}, nil
}

func (m *Money) Mul(n int) (*Money, error) {
	z, ok := mul(m.amount, n)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Money{
		amount:   z,
		currency: m.currency,
	}, nil
}

func (m *Money) Split(n int) ([]*Money, *Money, error) {
	if n <= 0 {
		return nil, nil, ErrSplitNegative
	}

	z, ok := div(m.amount, n)
	if !ok {
		return nil, nil, ErrInvalidOperation
	}

	ms := make([]*Money, n)
	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: z, currency: m.currency}
	}
	r, ok := mod(m.amount, n)
	if !ok {
		return nil, nil, ErrInvalidOperation
	}
	return ms, &Money{amount: r, currency: m.currency}, nil
}

func (m *Money) Alloc(rs ...int) ([]*Money, *Money, error) {
	if len(rs) == 0 {
		return nil, nil, ErrNoRatioSpecified
	}

	// sum of ratios
	var sum int
	for _, r := range rs {
		sum += r
	}

	var total amount
	var ms []*Money

	for _, r := range rs {
		a, ok := alloc(m.amount, r, sum)
		if !ok {
			return nil, nil, ErrInvalidOperation
		}
		party := &Money{
			amount:   a,
			currency: m.currency,
		}
		ms = append(ms, party)
		total, ok = add(total, m.amount)
		if !ok {
			return nil, nil, ErrInvalidOperation
		}
	}

	// leftover
	lo, ok := sub(m.amount, total)
	if !ok {
		return nil, nil, ErrInvalidOperation
	}
	return ms, &Money{amount: lo, currency: m.currency}, nil
}

func (m *Money) Round(k Kind) (*Money, error) {
	s, i := k.Rounding(m.currency)
	r, ok := round(m.amount, s, i)
	if !ok {
		return nil, ErrInvalidOperation
	}
	return &Money{amount: r, currency: m.currency}, nil
}

func (m *Money) Display() string {
	panic("not implemented")
}

func (m *Money) String() string {
	panic("not implemented")
}

func (m *Money) MarshalJSON() ([]byte, error) {
	panic("not implemented")
}
func (m *Money) UnmarshalJSON([]byte) error {
	panic("not implemented")
}

func (m *Money) Equals(om *Money) bool {
	return m.amount == om.amount &&
		m.currency.Equals(om.currency)
}

func (m *Money) GreaterThan(om *Money) (bool, error) {
	panic("not implemented")
}

func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	panic("not implemented")
}

func (m *Money) LessThan(om *Money) (bool, error) {
	panic("not implemented")
}

func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	panic("not implemented")
}

func (m *Money) IsPositive() bool {
	return m.amount >= 0
}

func (m *Money) IsZero() bool {
	return m.amount == 0
}

func (m *Money) IsNegative() bool {
	return m.amount < 0
}

func (m *Money) Abs() (bool, error) {
	panic("not implemented")
}

func (m *Money) Neg() (bool, error) {
	panic("not implemented")
}

//
// PRIVATE
//

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.currency.Equals(om.currency) {
		return ErrCurrencyMismatch
	}
	return nil
}
