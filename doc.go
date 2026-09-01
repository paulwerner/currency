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
// increments) and per-locale formatting data (number symbols, currency
// format patterns, currency symbols) are generated from the Unicode CLDR
// project; see internal/cldrgen and the Makefile for how to regenerate
// the tables.
//
// A minimal example:
//
//	price := currency.NewAmount(1999, currency.EUR) // 19.99 EUR
//	total, _ := price.Mul(3)
//	parts, rest, _ := total.Alloc(50, 30, 20)
//	fmt.Println(total, parts, rest)
//
// Amount.String renders a locale-independent "EUR 12.34" form;
// Amount.Display renders an amount for a BCP 47 locale using the CLDR
// data, e.g. "19,99 €" for "de" and "€19.99" for "en"; and
// Amount.DisplayKind adds kind-based display on top: the locale's
// accounting pattern for negative amounts ("($12.34)" for "en") and
// cash rounding at the cash scale ("CHF 10.15", "12 kr" for SEK).
package currency

//go:generate go run ./internal/cldrgen
