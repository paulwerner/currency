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
