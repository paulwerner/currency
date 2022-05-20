package money

import "math/big"

// constants for currency codes according to ISO 4217
const (
	EUR Currency = "EUR"
)

// types allowed for monetary unit
type (
	Int    = int64
	BigInt = big.Int
)
