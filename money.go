package money

import "math/big"

// types allowed for monetary unit
type (
	Int    = int64
	BigInt = big.Int
)

// Number represents a generic type considering int64 or big.Int defined in constants.go
type Number interface {
	Int | BigInt
}

// Amount represents the amount value used for computations
type Amount Number

// Money holds the monetary unit and currency
type Money[N Number] struct {
	amount   N
	currency *Currency
}

// New creates a new Money value for the given amount
// and currency code
func New[N Number](amount N, code string) *Money[N] {
	return &Money[N]{
		amount:   amount,
		currency: &Currency{code},
	}
}

// Currency returns the currency
func (m *Money[Number]) Currency() *Currency {
	return m.currency
}
