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
		t.Errorf("expected currency to be nil, got %s", currency)
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
		{newCurrency("USD"), newCurrency("USD"), true},
		{newCurrency("EUR"), newCurrency("EUR"), true},
		{newCurrency("EUR"), newCurrency("USD"), false},
		{newCurrency("USD"), newCurrency("EUR"), false},
		{newCurrency("USD"), newCurrency("EUR"), false},
	}
	for _, tc := range tcs {
		if ok := tc.c1.equals(tc.c2); ok != tc.want {
			t.Errorf("expected currencies %s and %s equality to be %v, got %v", tc.c1, tc.c2, ok, tc.want)
		}
	}
}
