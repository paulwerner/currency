package money

import "testing"

func TestGetCurrencyForCode(t *testing.T) {
	// test supported currency
	currency, err := GetCurrencyForCode("EUR")
	if err != nil {
		t.Errorf("Unexpected error occurred: %s", err)
	}
	if currency == nil || currency.code != EUR {
		t.Errorf("expected currency %s could not be found", EUR)
	}

	// test unsupported currency
	currency, err = GetCurrencyForCode("XXX")
	if err != ErrUnsupportedCurrency {
		t.Errorf("expected error %s, got %s", ErrUnsupportedCurrency, err)
	}
	if currency != nil {
		t.Errorf("expected currency to be nil, got %s", currency)
	}
}

func TestMustGetCurrencyForCode(t *testing.T) {
	// test must get currency
	currency := MustGetCurrencyForCode("EUR")
	if currency == nil || currency.code != EUR {
		t.Errorf("expected currency %s could not be found", EUR)
	}

	// test unsupported currency panics
	defer func() {
		if r := recover(); r != ErrUnsupportedCurrency {
			t.Errorf("expected error %s, got %s", ErrUnsupportedCurrency, r)
		}
	}()
	// panics
	MustGetCurrencyForCode("XXX")
}

func TestEquals(t *testing.T) {
	// test equal
	eur1 := MustGetCurrencyForCode("EUR")
	eur2 := MustGetCurrencyForCode("EUR")
	if eur1 != eur2 {
		t.Errorf("expected currencies %s and %s to be equal", eur1, eur2)
	}

	// test not equal
	usd := MustGetCurrencyForCode("USD")
	if eur1 == usd {
		t.Errorf("expected currencies %s and %s to not be equal", eur1, eur2)
	}
}
