package money

type Amount = int64

type Money struct {
	amount   Amount
	currency Currency
}

func (m *Money) Amount() Amount {
	return m.amount
}

func (m *Money) Currency() Currency {
	return m.currency
}

func (m *Money) Add(om *Money) (*Money, error) {
	panic("not implemented")
}

func (m *Money) Sub(om *Money) (*Money, error) {
	panic("not implemented")
}

func (m *Money) Div(d int64) (*Money, error) {
	panic("not implemented")
}

func (m *Money) Mul(mul int64) (*Money, error) {
	panic("not implemented")
}

func (m *Money) Round() (*Money, error) {
	panic("not implemented")
}

func (m *Money) Split(p int64) ([]*Money, error) {
	panic("not implemented")
}

func (m *Money) Alloc(ps ...int64) ([]*Money, error) {
	panic("not implemented")
}

func (m *Money) Display() string {
	panic("not implemented")
}
