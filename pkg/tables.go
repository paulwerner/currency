// Code generated. DO NOT EDIT.

package money

import "github.com/paulwerner/gomoney/internal/tag"

const (
	xxx = 288
	xts = 286
	usd = 253
	eur = 94
	jpy = 133
	gbp = 99
	chf = 61
	aud = 19
	nzd = 193
	cad = 58
	sek = 220
	nok = 191
	dkk = 82
	xag = 269
	xau = 270
	xpt = 283
	xpd = 281
	brl = 46
	cny = 68
	inr = 125
	rub = 211
	hkd = 114
	idr = 120
	krw = 141
	mxn = 179
	pln = 202
	sar = 214
	thb = 236
	try = 245
	twd = 247
	zar = 296
)

// currency holds an alphabetically sorted list of canonical 3-letter currency
// identifiers. Each identifier is followed by a byte of type currencyInfo,
// defined in gen_common.go.
const currency tag.Index = "" + // Size: 1220 bytes
	"\x00\x00\x00\x00\x41\x44\x50\x09\x41\x45\x44\x00\x41\x46\x41\x00\x41\x46" +
	"\x4e\x09\x41\x4c\x4b\x00\x41\x4c\x4c\x09\x41\x4d\x44\x08\x41\x4e\x47\x00" +
	"\x41\x4f\x41\x00\x41\x4f\x4b\x00\x41\x4f\x4e\x00\x41\x4f\x52\x00\x41\x52" +
	"\x41\x00\x41\x52\x4c\x00\x41\x52\x4d\x00\x41\x52\x50\x00\x41\x52\x53\x00" +
	"\x41\x54\x53\x00\x41\x55\x44\x00\x41\x57\x47\x00\x41\x5a\x4d\x00\x41\x5a" +
	"\x4e\x00\x42\x41\x44\x00\x42\x41\x4d\x00\x42\x41\x4e\x00\x42\x42\x44\x00" +
	"\x42\x44\x54\x00\x42\x45\x43\x00\x42\x45\x46\x00\x42\x45\x4c\x00\x42\x47" +
	"\x4c\x00\x42\x47\x4d\x00\x42\x47\x4e\x00\x42\x47\x4f\x00\x42\x48\x44\x1b" +
	"\x42\x49\x46\x09\x42\x4d\x44\x00\x42\x4e\x44\x00\x42\x4f\x42\x00\x42\x4f" +
	"\x4c\x00\x42\x4f\x50\x00\x42\x4f\x56\x00\x42\x52\x42\x00\x42\x52\x43\x00" +
	"\x42\x52\x45\x00\x42\x52\x4c\x00\x42\x52\x4e\x00\x42\x52\x52\x00\x42\x52" +
	"\x5a\x00\x42\x53\x44\x00\x42\x54\x4e\x00\x42\x55\x4b\x00\x42\x57\x50\x00" +
	"\x42\x59\x42\x00\x42\x59\x4e\x00\x42\x59\x52\x09\x42\x5a\x44\x00\x43\x41" +
	"\x44\x28\x43\x44\x46\x00\x43\x48\x45\x00\x43\x48\x46\x28\x43\x48\x57\x00" +
	"\x43\x4c\x45\x00\x43\x4c\x46\x24\x43\x4c\x50\x09\x43\x4e\x48\x00\x43\x4e" +
	"\x58\x00\x43\x4e\x59\x00\x43\x4f\x50\x08\x43\x4f\x55\x00\x43\x52\x43\x08" +
	"\x43\x53\x44\x00\x43\x53\x4b\x00\x43\x55\x43\x00\x43\x55\x50\x00\x43\x56" +
	"\x45\x00\x43\x59\x50\x00\x43\x5a\x4b\x08\x44\x44\x4d\x00\x44\x45\x4d\x00" +
	"\x44\x4a\x46\x09\x44\x4b\x4b\x30\x44\x4f\x50\x00\x44\x5a\x44\x00\x45\x43" +
	"\x53\x00\x45\x43\x56\x00\x45\x45\x4b\x00\x45\x47\x50\x00\x45\x52\x4e\x00" +
	"\x45\x53\x41\x00\x45\x53\x42\x00\x45\x53\x50\x09\x45\x54\x42\x00\x45\x55" +
	"\x52\x00\x46\x49\x4d\x00\x46\x4a\x44\x00\x46\x4b\x50\x00\x46\x52\x46\x00" +
	"\x47\x42\x50\x00\x47\x45\x4b\x00\x47\x45\x4c\x00\x47\x48\x43\x00\x47\x48" +
	"\x53\x00\x47\x49\x50\x00\x47\x4d\x44\x00\x47\x4e\x46\x09\x47\x4e\x53\x00" +
	"\x47\x51\x45\x00\x47\x52\x44\x00\x47\x54\x51\x00\x47\x57\x45\x00\x47\x57" +
	"\x50\x00\x47\x59\x44\x08\x48\x4b\x44\x00\x48\x4e\x4c\x00\x48\x52\x44\x00" +
	"\x48\x52\x4b\x00\x48\x54\x47\x00\x48\x55\x46\x08\x49\x44\x52\x08\x49\x45" +
	"\x50\x00\x49\x4c\x50\x00\x49\x4c\x52\x00\x49\x4c\x53\x00\x49\x4e\x52\x00" +
	"\x49\x51\x44\x09\x49\x52\x52\x09\x49\x53\x4a\x00\x49\x53\x4b\x09\x49\x54" +
	"\x4c\x09\x4a\x4d\x44\x00\x4a\x4f\x44\x1b\x4a\x50\x59\x09\x4b\x45\x53\x00" +
	"\x4b\x47\x53\x00\x4b\x48\x52\x00\x4b\x4d\x46\x09\x4b\x50\x57\x09\x4b\x52" +
	"\x48\x00\x4b\x52\x4f\x00\x4b\x52\x57\x09\x4b\x57\x44\x1b\x4b\x59\x44\x00" +
	"\x4b\x5a\x54\x00\x4c\x41\x4b\x09\x4c\x42\x50\x09\x4c\x4b\x52\x00\x4c\x52" +
	"\x44\x00\x4c\x53\x4c\x00\x4c\x54\x4c\x00\x4c\x54\x54\x00\x4c\x55\x43\x00" +
	"\x4c\x55\x46\x09\x4c\x55\x4c\x00\x4c\x56\x4c\x00\x4c\x56\x52\x00\x4c\x59" +
	"\x44\x1b\x4d\x41\x44\x00\x4d\x41\x46\x00\x4d\x43\x46\x00\x4d\x44\x43\x00" +
	"\x4d\x44\x4c\x00\x4d\x47\x41\x09\x4d\x47\x46\x09\x4d\x4b\x44\x00\x4d\x4b" +
	"\x4e\x00\x4d\x4c\x46\x00\x4d\x4d\x4b\x09\x4d\x4e\x54\x08\x4d\x4f\x50\x00" +
	"\x4d\x52\x4f\x09\x4d\x52\x55\x00\x4d\x54\x4c\x00\x4d\x54\x50\x00\x4d\x55" +
	"\x52\x08\x4d\x56\x50\x00\x4d\x56\x52\x00\x4d\x57\x4b\x00\x4d\x58\x4e\x00" +
	"\x4d\x58\x50\x00\x4d\x58\x56\x00\x4d\x59\x52\x00\x4d\x5a\x45\x00\x4d\x5a" +
	"\x4d\x00\x4d\x5a\x4e\x00\x4e\x41\x44\x00\x4e\x47\x4e\x00\x4e\x49\x43\x00" +
	"\x4e\x49\x4f\x00\x4e\x4c\x47\x00\x4e\x4f\x4b\x08\x4e\x50\x52\x00\x4e\x5a" +
	"\x44\x00\x4f\x4d\x52\x1b\x50\x41\x42\x00\x50\x45\x49\x00\x50\x45\x4e\x00" +
	"\x50\x45\x53\x00\x50\x47\x4b\x00\x50\x48\x50\x00\x50\x4b\x52\x08\x50\x4c" +
	"\x4e\x00\x50\x4c\x5a\x00\x50\x54\x45\x00\x50\x59\x47\x09\x51\x41\x52\x00" +
	"\x52\x48\x44\x00\x52\x4f\x4c\x00\x52\x4f\x4e\x00\x52\x53\x44\x09\x52\x55" +
	"\x42\x00\x52\x55\x52\x00\x52\x57\x46\x09\x53\x41\x52\x00\x53\x42\x44\x00" +
	"\x53\x43\x52\x00\x53\x44\x44\x00\x53\x44\x47\x00\x53\x44\x50\x00\x53\x45" +
	"\x4b\x08\x53\x47\x44\x00\x53\x48\x50\x00\x53\x49\x54\x00\x53\x4b\x4b\x00" +
	"\x53\x4c\x4c\x09\x53\x4f\x53\x09\x53\x52\x44\x00\x53\x52\x47\x00\x53\x53" +
	"\x50\x00\x53\x54\x44\x09\x53\x54\x4e\x00\x53\x55\x52\x00\x53\x56\x43\x00" +
	"\x53\x59\x50\x09\x53\x5a\x4c\x00\x54\x48\x42\x00\x54\x4a\x52\x00\x54\x4a" +
	"\x53\x00\x54\x4d\x4d\x09\x54\x4d\x54\x00\x54\x4e\x44\x1b\x54\x4f\x50\x00" +
	"\x54\x50\x45\x00\x54\x52\x4c\x09\x54\x52\x59\x00\x54\x54\x44\x00\x54\x57" +
	"\x44\x08\x54\x5a\x53\x08\x55\x41\x48\x00\x55\x41\x4b\x00\x55\x47\x53\x00" +
	"\x55\x47\x58\x09\x55\x53\x44\x00\x55\x53\x4e\x00\x55\x53\x53\x00\x55\x59" +
	"\x49\x09\x55\x59\x50\x00\x55\x59\x55\x00\x55\x59\x57\x24\x55\x5a\x53\x08" +
	"\x56\x45\x42\x00\x56\x45\x46\x08\x56\x45\x53\x00\x56\x4e\x44\x09\x56\x4e" +
	"\x4e\x00\x56\x55\x56\x09\x57\x53\x54\x00\x58\x41\x46\x09\x58\x41\x47\x00" +
	"\x58\x41\x55\x00\x58\x42\x41\x00\x58\x42\x42\x00\x58\x42\x43\x00\x58\x42" +
	"\x44\x00\x58\x43\x44\x00\x58\x44\x52\x00\x58\x45\x55\x00\x58\x46\x4f\x00" +
	"\x58\x46\x55\x00\x58\x4f\x46\x09\x58\x50\x44\x00\x58\x50\x46\x09\x58\x50" +
	"\x54\x00\x58\x52\x45\x00\x58\x53\x55\x00\x58\x54\x53\x00\x58\x55\x41\x00" +
	"\x58\x58\x58\x00\x59\x44\x44\x00\x59\x45\x52\x09\x59\x55\x44\x00\x59\x55" +
	"\x4d\x00\x59\x55\x4e\x00\x59\x55\x52\x00\x5a\x41\x4c\x00\x5a\x41\x52\x00" +
	"\x5a\x4d\x4b\x09\x5a\x4d\x57\x00\x5a\x52\x4e\x00\x5a\x52\x5a\x00\x5a\x57" +
	"\x44\x09\x5a\x57\x4c\x00\x5a\x57\x52\x00\xff\xff\xff\xff"

const numCurrencies = 303
