# GoMoney

GoMoney brings handling of monetary value to the go programming language. As the phrase says "time is money", it's build from the ground up for high performance and accuracy using up to date data from the Unicode CLDR Project.


## Quick start
`go get github.com/paulwerner/gomoney`


## Features
- low ops and memory allocation
- precise arithmetic (add, sub, split, multiply) with proper rounding
- support for 25+ currencies
- display format based on region or accounting format
- allow custom de-/serializer
- region and currency generation based on CLDR data
...


## Usage
...


## COde Generation
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
