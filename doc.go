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
// increments) is generated from the Unicode CLDR project; see
// internal/cldrgen and the Makefile for how to regenerate the tables.
//
// A minimal example:
//
//	price := currency.NewAmount(1999, currency.EUR) // 19.99 EUR
//	total, _ := price.Mul(3)
//	parts, rest, _ := total.Alloc(50, 30, 20)
//	fmt.Println(total, parts, rest)
//
// Locale-based formatting is planned but not yet implemented; see
// docs/ROADMAP.md for the current state of the roadmap.
package currency

//go:generate go run ./internal/cldrgen
