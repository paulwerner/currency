package money

import "testing"

func TestCurrency_GetCurrency(t *testing.T) {
	// test supported currency
	currency, err := GetCurrency("EUR")
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if currency == nil || currency.code != EUR {
		t.Errorf("expected currency %s could not be found", EUR)
	}

	// test unsupported currency
	currency, err = GetCurrency("XXX")
	if err != ErrUnsupportedCurrency {
		t.Errorf("expected error %s, got %s", ErrUnsupportedCurrency, err)
	}
	if currency != nil {
		t.Errorf("expected currency to be nil, got %v", currency)
	}
}

func TestCurrency_MustGetCurrency(t *testing.T) {
	// test must get currency
	currency := MustGetCurrency("EUR")
	if currency == nil || currency.code != EUR {
		t.Errorf("expected currency %s could not be found", EUR)
	}

	// test unsupported currency panics
	defer func() {
		if err := recover(); err != ErrUnsupportedCurrency {
			t.Errorf("expected error %s, got %s", ErrUnsupportedCurrency, err)
		}
	}()
	// panics
	MustGetCurrency("XXX")
}

func TestCurrency_Equal(t *testing.T) {
	tcs := []struct {
		c1   *Currency
		c2   *Currency
		want bool
	}{
		{newCurrency("USD", 2), newCurrency("USD", 2), true},
		{newCurrency("EUR", 2), newCurrency("EUR", 2), true},
		{newCurrency("EUR", 2), newCurrency("USD", 2), false},
		{newCurrency("USD", 2), newCurrency("EUR", 2), false},
		{newCurrency("USD", 2), newCurrency("EUR", 2), false},
	}
	for _, tc := range tcs {
		if ok := tc.c1.equals(tc.c2); ok != tc.want {
			t.Errorf("expected currencies %v and %v equality to be %v, got %v", tc.c1, tc.c2, ok, tc.want)
		}
	}
}
