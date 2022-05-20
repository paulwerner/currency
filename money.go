package money

// Number represents a generic type considering int64 or big.Int defined in constants.go
type Number interface {
	Int | BigInt
}

// Amount hold the amounts value used for computations
type Amount[N Number] struct {
	val N
}

// Money holds the monetary unit and currency
type Money[N Number] struct {
	amount   Amount[N]
	currency Currency
}

// New creates a new Money value for th given amount and currency
func New[N Number](amount N, currency Currency) *Money[N] {
	return &Money[N]{
		amount:   Amount[N]{amount},
		currency: currency,
	}
}
