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


## Development
Model:
======
Exposed:
-------

Money: struct
    - Amount : Amount
    - Currency : Currency

    Amount() : Amount
    Currency() : Currency
    Add(:Money) : Money
    Sub(:Money) : Money
    Div(:int) : Money
    Mul(:int) : Money
    Round() : Money
    Split(:int) : []Money
    Alloc(:...int) : []Money
    Display() : string
    String() : string

Amount : int64

Currency : struct 
    - index : uint16

    Code() : string
    CodeNumeric() : int
    String() : string

Internal:
--------

Kind : struct
    - rounding : Rounding
    
    Rounding(cu Currency) : Rounding
    Formatting(r Region) : Formatting

Rounding : struct
    - scale : int
    - increment : int

Formatting : struct 
    - region : Region
    - format : string

Region : 
    - language : string(ISO-639)



Tables:
const:
    currencyCode -> currencies.index

tag:

currency : struct
    - code string
    - {scale: int, fraction: int} : struct

regions : map[string][region]
    - language : string(ISO-639)
    - currencies : []currency

formats
