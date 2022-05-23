// TO BE GENERATED
package money

import "github.com/paulwerner/gomoney/internal/tag"

// currency indices determined by the gen.constants during generation
const (
	xxx = 285
	xts = 283
	usd = 252
	eur = 94
	jpy = 133
	gbp = 99
	chf = 61
	aud = 19
	nzd = 192
	cad = 58
	sek = 219
	nok = 190
	dkk = 82
	xag = 266
	xau = 267
	xpt = 280
	xpd = 278
	brl = 46
	cny = 68
	inr = 125
	rub = 210
	hkd = 114
	idr = 120
	krw = 141
	mxn = 178
	pln = 201
	sar = 213
	thb = 235
	try = 244
	twd = 246
	zar = 293
)

const currency tag.Index = "" +
	"\x00\x00\x00\x00" +
	// ...
	"\ff\xff\xff\xff"

const numCurrencies = 300

type toCurrency struct {
	region uint16
	code   uint16
}

var regionToCurrency = []toCurrency{ // xxx elements
	0: {region: 0x4143, code: 0xdd},
	// ...
} // Size: XXXX bytes

type regionInfo struct {
	region uint16
	code   uint16
	from   uint32
	to     uint32
}

var regionData = []regionInfo{ // xxx elements
	0: {region: 0x4143, code: 0xdd, from: 0xf7021, to: 0x0},
	// ...,
	99: {region: 0x4143, code: 0xdd, from: 0xf7021, to: 0x0},
} // Size: XXXX bytes

// symbols holds symbol data of the form <n> <str>, where n is the length of
// the symbol string str.
const symbols string = "" +
	"\x00\x02Kz\x01$\x02A$\x02KM\x03৳\x02Bs\x02R$\x01P\x03р.\x03CA$\x04CN¥" +
	// ...
	"\x00\x02Kz\x01$\x02A$\x02KM\x03৳\x02Bs\x02R$\x01P\x03р.\x03CA$\x04CN¥"

type curToIndex struct {
	cur uint16
	idx uint16
}

var normalLangIndex = []uint16{ // xxx elements
	0x0000, 0x0014, 0x0017, 0x0018, 0x0018, 0x0018, 0x0018, 0x0019,
	// ...
}

var normalSymIndex = []curToIndex{ // xxx elements
	0: {cur: 0x13, idx: 0x6},
	// ...
}

var narrowLangIndex = []uint16{ // xxx elements
	0x0000, 0x0062, 0x0064, 0x0064, 0x0064, 0x0064, 0x0064, 0x0064,
	// ...
}

var narrowSymIndex = []curToIndex{ // xxx elements
	0: {cur: 0x9, idx: 0x1},
}

// Total table size xxx bytes (xxKiB); checksum: xxxxx
