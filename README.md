# gomoney

gomoney provides a monetary value for the go programming language with precise operations, e.g. summing,  splitting, allocation and rounding, by using the currency's smallest fractional monetary unit for computation as defined in the ISO 4217.
For displaying the monetary value formats are based CLDR standard considering the selected region, as well as formats .

## Quick start
`go get github.com/paulwerner/gomoney`

## Features
- low ops and memory allocation
- precise arithmetic (add, sub, split, multiply) with proper rounding
- display format based on region or accounting format
- allow custom de-/serializer
- region and currency generation based on CLDR data
...
...

## Usage
### code generation
`go generate`

...


## Inspired by
- https://github.com/Rhymond/go-money
- https://pkg.go.dev/golang.org/x/text/currency

## Contribution
...