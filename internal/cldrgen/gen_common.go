// This file contains code common to the generator and the currency package
// itself. It is copied to the repository root as common.go by the generator
// (see gen.go), adjusting the package clause.
package main

import "time"

const (
	cashShift = 3
	roundMask = 0x7

	nonTenderBit = 0x8000
)

type currencyInfo byte

type roundingType struct {
	scale, increment uint8
}

var roundings = [...]roundingType{
	{2, 1}, // default
	{0, 1},
	{1, 1},
	{3, 1},
	{4, 1},
	{2, 5}, // cash rounding alternative
	{2, 50},
}

// localePattern is a CLDR currency format pattern, parsed at generation
// time into the literal affixes surrounding the formatted number. Affixes
// contain "¤" as the placeholder for the currency symbol and "-" as
// the placeholder for the locale's minus sign; everything else is literal
// text. primGroup and secGroup are the sizes of the rightmost and the
// subsequent integer digit groups ("#,##0.00" is 3/3, "#,##,##0.00" is
// 3/2); both are 0 when the pattern does not group digits.
type localePattern struct {
	posPrefix, posSuffix string
	negPrefix, negSuffix string
	primGroup, secGroup  uint8
}

// currencySymbol maps a 3-letter ISO 4217 code to the symbol a locale
// displays it with, e.g. "USD" -> "$".
type currencySymbol struct {
	code, symbol string
}

// localeData holds the latn number-formatting data of one locale as a
// sparse overlay over its parent locale: a zero-valued field ("" or 0)
// means the value is inherited from the parent chain. parent is the index
// of the parent locale in the locales table; the root locale is stored at
// index 0 and is its own parent. standard and accounting are 1-based
// indices into localePatterns (0 = inherited). symbols lists only the
// currency symbols this locale overrides, sorted by currency code.
type localeData struct {
	name    string // canonical BCP 47 tag, e.g. "de-CH"
	parent  uint16
	decimal string // decimal separator
	group   string // integer digit group separator
	minus   string // minus sign

	standard   uint16
	accounting uint16

	symbols []currencySymbol
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

func fromDate(date uint32) time.Time {
	return time.Date(
		int(date>>9),
		time.Month((date>>5)&0xf),
		int(date&0x1f), 0, 0, 0, 0, time.UTC,
	)
}
