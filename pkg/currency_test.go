package money

import (
	"testing"
)

func TestCurrency_Equals(t *testing.T) {
	tcs := []struct {
		c1   Currency
		c2   Currency
		want bool
	}{
		{EUR, EUR, true},
		{USD, EUR, false},
		{EUR, USD, false},
		{USD, USD, true},
	}

	for _, tc := range tcs {
		if eq := tc.c1.Equals(&tc.c2); eq != tc.want {
			t.Errorf("expected %v and %v to be %v, got %v", tc.c1, tc.c2, tc.want, eq)
		}
	}
}

func TestKind_Rounding(t *testing.T) {
	for i, tc := range []struct {
		kind          Kind
		cur           *Currency
		wantScale     int
		wantIncrement int
	}{
		{Standard, &USD, 2, 1},
		{Standard, &EUR, 2, 1},
		{Standard, &JPY, 0, 1},
		{Standard, &GBP, 2, 1},
		{Standard, &CHF, 2, 1},
		{Standard, &AUD, 2, 1},
		{Standard, &NZD, 2, 1},
		{Standard, &CAD, 2, 1},
		{Standard, &SEK, 2, 1},
		{Standard, &NOK, 2, 1},

		{Standard, &BRL, 2, 1},
		{Standard, &CNY, 2, 1},
		{Standard, &DKK, 2, 1},
		{Standard, &INR, 2, 1},
		{Standard, &RUB, 2, 1},
		{Standard, &HKD, 2, 1},
		{Standard, &IDR, 2, 1},
		{Standard, &KRW, 0, 1},
		{Standard, &MXN, 2, 1},
		{Standard, &PLN, 2, 1},
		{Standard, &SAR, 2, 1},
		{Standard, &THB, 2, 1},
		{Standard, &TRY, 2, 1},
		{Standard, &TWD, 2, 1},
		{Standard, &ZAR, 2, 1},

		// standard scale and increment for precious metals
		{Standard, &XAG, 2, 1},
		{Standard, &XAU, 2, 1},
		{Standard, &XPT, 2, 1},
		{Standard, &XPD, 2, 1},
	} {
		s, incr := tc.kind.Rounding(tc.cur)
		if s != tc.wantScale {
			t.Errorf("[%v]:[%s] want scale %v, got %v", i, tc.cur, tc.wantScale, s)
		}
		if incr != tc.wantIncrement {
			t.Errorf("[%v]:[%s] want increment %v, got %v", i, tc.cur, tc.wantIncrement, incr)
		}
	}
}

func TestCurrency_CurrencyFromISO(t *testing.T) {
	isos := []string{"XCD", "ALL", "ALK", "AMD", "RUR", "SUR", "AOA", "AOR", "AON"}

	for _, iso := range isos {
		cur, err := CurrencyFromISO(iso)
		if err != nil {
			t.Errorf("[%v]: expected error to be nil, got %v", iso, err)
		}
		if cur == nil {
			t.Errorf("[%v]: expected currency to be not nil, got %v", iso, err)
		}
		if cur.Code() != iso {
			t.Errorf("[%v]: error expected currency iso, got %v", iso, cur.Code())
		}
	}

	// not supported currency
	iso := "BTC"
	cur, err := CurrencyFromISO(iso)
	if err == nil {
		t.Errorf("[%v]: expected error not to be nil, got %v", iso, err)
	}
	if cur != nil {
		t.Errorf("[%v]: expected currency to be  nil, got %v", iso, err)
	}
}
