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
	panic("not implemented")
}

func toDate(t time.Time) uint32 {
	panic("not implemented")
}

func fromDate(date uint32) time.Time {
	panic("not implemented")
}
