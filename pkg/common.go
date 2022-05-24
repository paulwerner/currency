// to be generated from gen_common.go

package money

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

func regionToCode(r Region) uint16 {
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

func fromDate(date uint32) time.Time {
	return time.Date(int(date>>9), time.Month((date>>5)&0xf), int(date&0x1f), 0, 0, 0, 0, time.UTC)
}
