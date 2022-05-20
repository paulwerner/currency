package money

import "math/big"

// constants for currency codes according to ISO 4217
var (
	// Undefined and testing
	XXX Currency = Currency{}
	XTS Currency = Currency{"XTS"}

	// G10 currencies https://en.wikipedia.org/wiki/G10_currencies
	USD Currency = Currency{"USD"}
	EUR Currency = Currency{"EUR"}
	JPY Currency = Currency{"JPY"}
	GDP Currency = Currency{"GDP"}
	CHF Currency = Currency{"CHF"}
	AUD Currency = Currency{"AUD"}
	NZD Currency = Currency{"NZD"}
	CAD Currency = Currency{"CAD"}
	SEK Currency = Currency{"SEK"}
	NOK Currency = Currency{"NOK"}

	// Additional common currencies
	BRL Currency = Currency{"BRL"}
	CNY Currency = Currency{"CNY"}
	DKK Currency = Currency{"DKK"}
	INR Currency = Currency{"INR"}
	RUB Currency = Currency{"RUB"}
	HKD Currency = Currency{"HKD"}
	IDR Currency = Currency{"IDR"}
	KRW Currency = Currency{"KRW"}
	MXN Currency = Currency{"MXN"}
	PLN Currency = Currency{"PLN"}
	SAR Currency = Currency{"SAR"}
	THB Currency = Currency{"THB"}
	TRY Currency = Currency{"TRY"}
	TWD Currency = Currency{"TWD"}
	ZAR Currency = Currency{"ZAR"}

	// Precious metals
	XAG Currency = Currency{"XAG"}
	XAU Currency = Currency{"XAU"}
	XPT Currency = Currency{"XPT"}
	XPD Currency = Currency{"XPD"}
)

// types allowed for monetary unit
type (
	Int    = int64
	BigInt = big.Int
)
