// Package currency provides types and functions for handling monetary
// values safely and efficiently.
//
// Monetary values are represented by Amount, a pairing of an integer value
// in the currency's minor units (e.g. cents for USD, yen for JPY) with a
// Currency. All arithmetic (Add, Sub, Mul, Split, Alloc, Round) is
// overflow-checked and returns an error instead of silently wrapping
// around.
//
// Currency metadata (ISO 4217 codes, rounding scales and cash rounding
// increments) as well as locale data (currency symbols and number
// formatting patterns) are generated from the Unicode CLDR project; see
// gen.go and the Makefile for how to regenerate the tables.
//
// A minimal example:
//
//	price, _ := currency.NewAmount(1999, currency.EUR) // 19.99 EUR
//	tax, _ := price.Mul(19)
//	tax, _ = tax.Split(100) // ... or use Alloc for exact distribution
//
// Formatting is locale aware:
//
//	price.Display("de") // "19,99 €"
//	price.Display("en") // "€19.99"
package currency

//go:generate go run ./internal/cldrgen
