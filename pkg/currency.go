package money


type Currency struct {
	index uint16
}

func (c *Currency) Code() string {
	panic("not implemented")
}

func (c *Currency) CodeNumeric() uint8 {
	panic("not implemented")
}

func (c *Currency) String() string {
	panic("not implemented")
}

func (c *Currency) template(r Region) string {
	panic("not implemented")
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
