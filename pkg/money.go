package money

type Money struct {
	amount   *Amount
	currency *Currency
}

func New(v int64, cur *Currency) (*Money, error) {
	return &Money{
		amount:   amount(v),
		currency: cur,
	}, nil
}

func NewFromISO(v int64, iso string) (*Money, error) {
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
	panic("not implemented")
}

func (m *Money) Amount() int64 {
	panic("not implemented")
}

func (m *Money) SameCurrency(om *Money) bool {
	panic("not implemented")
}

func (m *Money) Add(om *Money) (*Money, error) {
	panic("not implemented")
}

func (m *Money) Sub(om *Money) (*Money, error) {
	panic("not implemented")
}

func (m *Money) Mul(mul int64) (*Money, error) {
	panic("not implemented")
}

func (m *Money) Split(d int64) ([]*Money, error) {
	panic("not implemented")
}

func (m *Money) SplitWithRemainder(d int64) ([]*Money, *Money, error) {
	panic("not implemented")
}

func (m *Money) Alloc(r int64, s *Money) ([]*Money, error) {
	panic("not implemented")
}

func (m *Money) AllocWithRemainder(r int64, s *Money) ([]*Money, *Money, error) {
	panic("not implemented")
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
