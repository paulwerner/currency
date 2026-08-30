# Currency

[![CI](https://github.com/paulwerner/currency/actions/workflows/ci.yml/badge.svg)](https://github.com/paulwerner/currency/actions/workflows/ci.yml)

Currency is a library for handling monetary values in Go (Golang).
Amounts are integer values in a currency's minor units (cents, yen, …)
with overflow-checked arithmetic, and the currency metadata — ISO 4217
codes, rounding scales, cash rounding increments — is generated from the
Unicode [CLDR](https://cldr.unicode.org/) project.

## Quick start

```sh
go get github.com/paulwerner/currency
```

```go
package main

import (
	"fmt"

	"github.com/paulwerner/currency"
)

func main() {
	price := currency.NewAmount(1999, currency.EUR) // 19.99 EUR
	total, _ := price.Mul(3)

	// distribute 59.97 EUR by ratios 50/30/20, nothing gets lost
	parts, rest, _ := total.Alloc(50, 30, 20)
	fmt.Println(total, parts[0], rest) // EUR 59.97 EUR 29.98 EUR 0.01

	// cash rounding: 0.05 CHF increments
	cash := currency.NewAmount(1002, currency.CHF)
	rounded, _ := cash.Round(currency.Cash)
	fmt.Println(rounded) // CHF 10.00

	// locale-based formatting
	fmt.Println(price.Display("de")) // 19,99 €
	fmt.Println(price.Display("en")) // €19.99

	// kind-based display: accounting parentheses, cash rounding
	debt, _ := price.Neg()
	fmt.Println(debt.DisplayKind("en", currency.Accounting)) // (€19.99)
	fmt.Println(cash.DisplayKind("en", currency.Cash))       // CHF 10.00

	// parsing and (de-)serialization
	a, _ := currency.NewFromISO(500, "DKK")
	b, _ := a.MarshalJSON()
	fmt.Println(string(b)) // {"amount":500,"currency":"DKK"}
}
```

## Features

- integer minor-unit representation, low ops and memory allocation
- precise arithmetic with overflow checks: `Add`, `Sub`, `Mul`, `Split`,
  `Alloc` (ratio allocation without losing a cent), `Abs`, `Neg`
- comparisons with currency-mismatch detection: `Cmp`, `GreaterThan`,
  `LessThan`, …
- kind-based rounding: `Standard`, `Cash` (e.g. 0.05 CHF, 0.50 DKK),
  `Accounting`
- support for 300+ currencies, data generated from the CLDR project
- JSON de-/serialization
- locale-based formatting from CLDR data: `Display("de-CH")` → `EUR 1’234.56`,
  with BCP 47 tag normalization and parent-chain fallback
- kind-based displaying: `DisplayKind("en", currency.Accounting)` →
  `($12.34)` for negatives, `DisplayKind("sv", currency.Cash)` → `12 kr`
  (cash rounding at the cash scale)

## Development

See [AGENTS.md](AGENTS.md) for layout, conventions, and how to build and
test, and [docs/ROADMAP.md](docs/ROADMAP.md) for the roadmap's design
notes and the change log.

```sh
make build test vet fmt-check   # what CI runs
make gen-fetch                  # refetch CLDR data and regenerate tables
```

## Inspired by

- https://github.com/Rhymond/go-money
- https://pkg.go.dev/golang.org/x/text/currency
