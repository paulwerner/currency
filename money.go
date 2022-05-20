package money

import "math/big"

// types allowed for monetary unit
type (
	Int    = int64
	BigInt = big.Int
)

// Number represents a generic type considering int64 or big.Int defined in constants.go
type Amount interface {
	Int | BigInt
}

// Money holds the monetary unit and currency
type Money[A Amount] struct {
	amount   A
	currency *Currency
}

// New creates a new Money value for the given amount
// and currency code
func New[A Amount](amount A, code string) *Money[A] {
	return &Money[A]{
		amount:   amount,
		currency: &Currency{code},
	}
}

// Currency returns the currency
func (m *Money[Number]) Currency() *Currency {
	return m.currency
}
