package money

// Kind determines the rounding and rendering properties of the currency value
type Kind struct {
	rounding rounding
	// TODO: formatting: (standard|accounting)
}

type rounding byte

const (
	standard rounding = iota
	cash
)

// Rounding reports the rounding characteristics for the given currency, where
// scale is the number of fractional decimals and increment is the number of
// units in terms of 10^(-scale) to which to round to
func (k Kind) Rounding(cur Currency) (scale, increment int) {
	info := currency.Elem(int(cur.index))[3]
	switch k.rounding {
	case standard:
		info &= roundMask
	case cash:
		info >>= cashShift
	}
	return int(roundings[info].scale), int(roundings[info].increment)
}

type Currency struct {
	index uint16
}

// Code reports the currency's ISO code
// See: Currency.String()
func (c *Currency) Code() string {
	return c.String()
}

// String returns the currency's ISO code
func (c *Currency) String() string {
	if c.index == 0 {
		return "XXX"
	}
	return currency.Elem(int(c.index))[:3]
}

var (
	// Undefined and testing
	XXX Currency = Currency{}
	XTS Currency = Currency{xts}

	// G10 currencies https://en.wikipedia.org/wiki/G10_currencies
	USD Currency = Currency{usd}
	EUR Currency = Currency{eur}
	JPY Currency = Currency{jpy}
	GBP Currency = Currency{gbp}
	CHF Currency = Currency{chf}
	AUD Currency = Currency{aud}
	NZD Currency = Currency{nzd}
	CAD Currency = Currency{cad}
	SEK Currency = Currency{sek}
	NOK Currency = Currency{nok}

	// Additional common currencies as defined by CLDR
	BRL Currency = Currency{brl}
	CNY Currency = Currency{cny}
	DKK Currency = Currency{dkk}
	INR Currency = Currency{inr}
	RUB Currency = Currency{rub}
	HKD Currency = Currency{hkd}
	IDR Currency = Currency{idr}
	KRW Currency = Currency{krw}
	MXN Currency = Currency{mxn}
	PLN Currency = Currency{pln}
	SAR Currency = Currency{sar}
	THB Currency = Currency{thb}
	TRY Currency = Currency{try}
	TWD Currency = Currency{twd}
	ZAR Currency = Currency{zar}

	// Precious metals
	XAG Currency = Currency{xag}
	XAU Currency = Currency{xau}
	XPT Currency = Currency{xpt}
	XPD Currency = Currency{xpd}
)
