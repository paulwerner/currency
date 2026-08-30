package currency

import (
	"math"
	"strconv"
	"testing"
)

// Expected strings below spell out the exact separators:   is a
// no-break space,   a narrow no-break space, − (U+2212) a real
// minus sign, ’ (U+2019) the Swiss apostrophe group separator
func TestAmount_Display(t *testing.T) {
	for _, tc := range []struct {
		value  int
		cur    Currency
		locale string
		want   string
	}{
		// the roadmap examples
		{1999, EUR, "en", "€19.99"},
		{1999, EUR, "de", "19,99 €"},

		{-1999, EUR, "en", "-€19.99"},
		{-1999, EUR, "de", "-19,99 €"},
		{0, USD, "en", "$0.00"},

		// de-CH: own separators and pattern, EUR symbol reset to the
		// ISO code; the negative pattern puts the minus after the symbol
		{123456, EUR, "de-CH", "EUR 1’234.56"},
		{-123456, EUR, "de-CH", "EUR-1’234.56"},

		// currency spacing: en has no CHF symbol, the ISO-code fallback
		// ends in a letter and gets a no-break space before the digits
		{1234, CHF, "en", "CHF 12.34"},
		{-1234, CHF, "en", "-CHF 12.34"},

		// en-IN uses Indian 3/2 digit grouping
		{123456789, USD, "en-IN", "$12,34,567.89"},

		// fr: narrow no-break space grouping, trailing symbol
		{123456, USD, "fr", "1 234,56 $US"},
		// sv: U+2212 minus sign, no-break space grouping
		{-123456, SEK, "sv", "−1 234,56 kr"},
		// ru: kopeck-scaled ruble sign
		{-123456, RUB, "ru", "-1 234,56 ₽"},
		// nb: inherits the "kr" symbol from the flattened "no" parent
		{123456, NOK, "nb", "kr 1 234,56"},

		// scale-0 currency: no decimal part at all
		{1234567, JPY, "ja", "￥1,234,567"},
		{1234567, JPY, "de", "1.234.567 ¥"},
		{123456, JPY, "en-GB", "JP¥123,456"},

		// es-MX overrides the es separators back to "."/"," and shows
		// the peso as a plain "$"
		{12345, MXN, "es-MX", "$123.45"},
		{1234567, EUR, "es", "12.345,67 €"},
		{-1234, TWD, "zh-Hant", "-$12.34"},
		// RTL data is used as-is; no bidi isolation marks are emitted
		{1234, SAR, "ar", "ر.س.‏ 12.34"},

		// unknown locales fall back to root, unknown symbols to the
		// ISO code, XXX to its CLDR placeholder symbol
		{1234, CHF, "klingon", "CHF 12.34"},
		{1234, USD, "", "US$ 12.34"},
		{1234, XXX, "en", "¤12.34"},
	} {
		a := NewAmount(tc.value, tc.cur)
		if got := a.Display(tc.locale); got != tc.want {
			t.Errorf("[%v %v %v]: want %q, got %q", tc.locale, tc.cur, tc.value, tc.want, got)
		}
	}
}

// TestAmount_Display_LocaleLookup checks tag normalization ("_" vs "-",
// case folding) and the truncating fallback chain
func TestAmount_Display_LocaleLookup(t *testing.T) {
	for _, tc := range []struct {
		value  int
		cur    Currency
		locale string
		want   string
	}{
		// separator and case variants of de-CH
		{123456, EUR, "de_CH", "EUR 1’234.56"},
		{123456, EUR, "DE-ch", "EUR 1’234.56"},
		{123456, EUR, "de_ch", "EUR 1’234.56"},
		// script subtags are title-cased
		{-1234, TWD, "ZH-HANT", "-$12.34"},
		// unknown region truncates to the language
		{1234, USD, "en-US", "$12.34"},
		{1999, EUR, "de-DE", "19,99 €"},
		// multi-step truncation: zh-Hant-TW → zh-Hant
		{-1234, TWD, "zh-hant-tw", "-$12.34"},
		{-1234, TWD, "zh_Hant_TW", "-$12.34"},
		// documented limitation: no likely-subtags mapping, so zh-TW
		// truncates to zh (Simplified), which shows TWD via root
		{-1234, TWD, "zh-TW", "-NT$12.34"},
		// region locales that do exist are used directly
		{1234, CAD, "fr-CA", "12,34 $"},
		{1234, AUD, "en-AU", "$12.34"},
		// explicit root and a fully unknown tag behave alike
		{1234, USD, "root", "US$ 12.34"},
		{1234, USD, "xx-YY", "US$ 12.34"},
	} {
		a := NewAmount(tc.value, tc.cur)
		if got := a.Display(tc.locale); got != tc.want {
			t.Errorf("[%v %v %v]: want %q, got %q", tc.locale, tc.cur, tc.value, tc.want, got)
		}
	}
}

// TestAmount_Display_Bounds pins the loBound behavior: the sign is
// carried by the pattern, so the digits must not be negated (which
// would overflow for MinInt in a scale-0 currency)
func TestAmount_Display_Bounds(t *testing.T) {
	wantJPY := map[int]string{
		32: "¥2,147,483,648",
		64: "¥9,223,372,036,854,775,808",
	}[strconv.IntSize]
	if got := NewAmount(math.MinInt, JPY).Display("en"); got != "-"+wantJPY {
		t.Errorf("[JPY]: want %q, got %q", "-"+wantJPY, got)
	}

	wantEUR := map[int]string{
		32: "€21,474,836.48",
		64: "€92,233,720,368,547,758.08",
	}[strconv.IntSize]
	if got := NewAmount(math.MinInt, EUR).Display("en"); got != "-"+wantEUR {
		t.Errorf("[EUR]: want %q, got %q", "-"+wantEUR, got)
	}
	if got := NewAmount(math.MaxInt, EUR).Display("en"); got != wantEUR[:len(wantEUR)-1]+"7" {
		t.Errorf("[EUR max]: want %q, got %q", wantEUR[:len(wantEUR)-1]+"7", got)
	}
}

func TestGroupDigits(t *testing.T) {
	for _, tc := range []struct {
		digits    string
		sep       string
		prim, sec uint8
		want      string
	}{
		{"0", ",", 3, 3, "0"},
		{"999", ",", 3, 3, "999"},
		{"1000", ",", 3, 3, "1,000"},
		{"123456", ",", 3, 3, "123,456"},
		{"1234567", ",", 3, 3, "1,234,567"},
		{"9223372036854775808", ",", 3, 3, "9,223,372,036,854,775,808"},
		// Indian 3/2 grouping
		{"1234", ",", 3, 2, "1,234"},
		{"12345", ",", 3, 2, "12,345"},
		{"123456", ",", 3, 2, "1,23,456"},
		{"1234567", ",", 3, 2, "12,34,567"},
		{"123456789", ",", 3, 2, "12,34,56,789"},
		// prim 0 disables grouping; sec 0 falls back to prim
		{"1234567", ",", 0, 0, "1234567"},
		{"1234567", ",", 3, 0, "1,234,567"},
		// multi-byte separator
		{"1234567", " ", 3, 3, "1 234 567"},
	} {
		if got := groupDigits(tc.digits, tc.sep, tc.prim, tc.sec); got != tc.want {
			t.Errorf("[%v %d/%d]: want %q, got %q", tc.digits, tc.prim, tc.sec, tc.want, got)
		}
	}
}
