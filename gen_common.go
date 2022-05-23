//go:build ignore
// +build ignore

package main

import (
	"time"

	"golang.org/x/text/language"
)

// This file contains code common to gen.go and the package code

const (
	cashShift = 3
	roundMask = 0x7

	nonTenderBit = 0x8000
)

// currencyInfo contains the currency's information
// bits 0..2: index into roundings for standard rounding
// bits 3..5: index into roundings for cash rounding
type currencyInfo byte

// roundingType defines the scale (number of fractional decimals) and increments
// in term of units of size 10^-scale, e.g. scale == 2 and increment == 1,
// the currency is rounded to units of 0.01.
type roundingType struct {
	scale, increment uint8
}

// roundings contains rounding data for currencies
var roundings = [...]roundingType{
	{2, 1}, // default
	{0, 1},
	{1, 1},
	{3, 1},
	{4, 1},
	{2, 5}, // cash rounding alternative
	{2, 50},
}

// regionToCode returns a 16-bit region code. Only two-letter codes are
// supported. (Three-letter codes are not needed.)
func regionToCode(r language.Region) uint16 {
	if s := r.String(); len(s) == 2 {
		return uint16(s[0])<<8 | uint16(s[1])
	}
	return 0
}

func toDate(t time.Time) uint32 {
	y := t.Year()
	if y == 1 {
		return 0
	}
	date := uint32(y) << 4
	date |= uint32(t.Month())
	date <<= 5
	date |= uint32(t.Day())
	return date
}
