# Currency

Currency is a library for handling monetary values in Go (Golang). 
Its data is based on the CLDR project which can be updated for newer versions using the internal generator.

## Quick start
`go get github.com/paulwerner/currency`


## Features
- low ops and memory allocation
- precise arithmetic with overflow checks
- support for 300+ currencies
- locale based displaying
- allow custom de-/serializer
- data generated based on the CLDR project


## Development
- [x] add calculator
- [x] add CLDR data fetching
- [x] add code generator
- [x] add currency data generation
- [ ] add locale data generation
- [ ] add locale based formatting
- [ ] add kind based displaying (standard, cash, accounting)
- [ ] add de-/serialization
- [ ] add documentation


## Code Generation
`make gen`


## Inspired By
- https://github.com/Rhymond/go-money
- https://pkg.go.dev/golang.org/x/text/currency

