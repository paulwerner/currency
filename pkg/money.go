package money

import "errors"

var (
	ErrCurrencyMismatch  = errors.New("money: currency mismatch")
	ErrOperationOverflow = errors.New("money: operation overflow")
	ErrSplitNegative     = errors.New("money: split must be positive")
	ErrNoRatioSpecified  = errors.New("money: no ratio specified")
)

type Money struct {
	amount   *Amount
	currency *Currency
}

func New(v int, cur *Currency) (*Money, error) {
	return &Money{
		amount:   amount(v),
		currency: cur,
	}, nil
}

func NewFromISO(v int, iso string) (*Money, error) {
	cur, err := CurrencyFromISO(iso)
	if err != nil {
		return nil, err
	}
	return &Money{
		amount:   amount(v),
		currency: cur,
	}, nil
}

func (m *Money) Currency() *Currency {
	return m.currency
}

func (m *Money) Amount() *Amount {
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
		return nil, ErrOperationOverflow
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
		return nil, ErrOperationOverflow
	}
	return &Money{
		amount:   z,
		currency: m.currency,
	}, nil
}

func (m *Money) Mul(n int) (*Money, error) {
	z, ok := mul(m.amount, n)
	if !ok {
		return nil, ErrOperationOverflow
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
		return nil, nil, ErrOperationOverflow
	}

	ms := make([]*Money, n)
	for i := 0; i < n; i++ {
		ms[i] = &Money{amount: z, currency: m.currency}
	}

	r := mod(m.amount, n)
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

	var total *Amount
	var ms []*Money

	for _, r := range rs {
		a, ok := alloc(m.amount, r, sum)
		if !ok {
			return nil, nil, ErrOperationOverflow
		}
		party := &Money{
			amount:   a,
			currency: m.currency,
		}
		ms = append(ms, party)
		total, ok = add(total, m.amount)
	}

	// leftover
	lo, ok := sub(m.amount, total)
	if !ok {
		return nil, nil, ErrOperationOverflow
	}
	return ms, &Money{amount: lo, currency: m.currency}, nil
}

func (m *Money) Round() *Money {
	panic("not implemented")
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
	panic("not implemented")
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

func (m *Money) IsPositive() (bool, error) {
	panic("not implemented")
}

func (m *Money) IsZero() (bool, error) {
	panic("not implemented")
}

func (m *Money) IsNegative() (bool, error) {
	panic("not implemented")
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
	panic("not implemented")
}

func (m *Money) compare(om *Money) int {
	panic("not implemented")
}
