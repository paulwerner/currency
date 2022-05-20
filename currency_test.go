package money

import "testing"

func TestGetCurrency(t *testing.T) {
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

func TestMustGetCurrency(t *testing.T) {
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

func TestEqual(t *testing.T) {
	// test equal
	eur1 := MustGetCurrency("EUR")
	eur2 := MustGetCurrency("EUR")
	if !eur1.Equal( eur2) {
		t.Errorf("expected currencies %s and %s to be equal", eur1, eur2)
	}

	// test not equal
	usd := MustGetCurrency("USD")
	if eur1.Equal(usd) {
		t.Errorf("expected currencies %s and %s to not be equal", eur1, usd)
	}
}
