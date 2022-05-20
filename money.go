package money

type Money struct {
	amount   int64
	currency *Currency
}

func New(amount int64, code string) *Money {
	return &Money{
		amount:   amount,
		currency: &Currency{code},
	}
}

func (m *Money) Currency() *Currency {
	return m.currency
}

func (m *Money) CurrencyCode() string {
	return m.currency.code
}

func (m *Money) Amount() int64 {
	return m.amount
}
